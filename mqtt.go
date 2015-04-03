package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"sync"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var MaxClientIdLen = 8
var MaxRetryCount = 3

type MQTTClient struct {
	Client     *MQTT.Client
	Opts       *MQTT.ClientOptions
	RetryCount int
	Subscribed map[string]byte

	lock *sync.Mutex // use for reconnect
}

// Connects connect to the MQTT broker with Options.
func (m *MQTTClient) Connect() (*MQTT.Client, error) {

	m.Client = MQTT.NewClient(m.Opts)

	log.Info("connecting...")

	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return m.Client, nil
}

func (m *MQTTClient) Publish(topic string, payload []byte, qos int, retain bool, sync bool) error {
	token := m.Client.Publish(topic, byte(qos), retain, payload)

	if sync == true {
		token.Wait()
	}

	return token.Error()
}

func (m *MQTTClient) Disconnect() error {
	if m.Client.IsConnected() {
		m.Client.Disconnect(20)
		log.Info("client disconnected")
	}
	return nil
}

func (m *MQTTClient) SubscribeOnConnect(client *MQTT.Client) {
	log.Infof("client connected")

	if len(m.Subscribed) > 0 {
		token := client.SubscribeMultiple(m.Subscribed, m.onMessageReceived)
		token.Wait()
		if token.Error() != nil {
			log.Error(token.Error())
		}
	}
}

func (m *MQTTClient) ConnectionLost(client *MQTT.Client, reason error) {
	log.Errorf("client disconnected: %s", reason)
}

func (m *MQTTClient) onMessageReceived(client *MQTT.Client, message MQTT.Message) {
	log.Infof("topic:%s / msg:%s", message.Topic(), message.Payload())
	fmt.Println(string(message.Payload()))
}

func getCertPool(pemPath string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()

	pemData, err := ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	return certs, nil
}

// getRandomClientId returns randomized ClientId.
func getRandomClientId() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, MaxClientIdLen)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return "mqttcli-" + string(bytes)
}

// NewOption returns ClientOptions via parsing command line options.
func NewOption(c *cli.Context) (*MQTT.ClientOptions, error) {
	opts := MQTT.NewClientOptions()

	host := c.String("host")
	port := c.Int("p")

	if host == "" {
		getSettingsFromFile(c.String("conf"), opts)
	}

	clientId := c.String("i")
	if clientId == "" {
		clientId = getRandomClientId()
	}
	opts.SetClientID(clientId)

	TLSConfig := &tls.Config{InsecureSkipVerify: false}
	cafile := c.String("cafile")
	scheme := "tcp"
	if cafile != "" {
		scheme = "ssl"
		certPool, err := getCertPool(cafile)
		if err != nil {
			return nil, err
		}
		TLSConfig.RootCAs = certPool
	}
	insecure := c.Bool("insecure")
	if insecure {
		TLSConfig.InsecureSkipVerify = true
	}
	opts.SetTLSConfig(TLSConfig)

	user := c.String("u")
	if user != "" {
		opts.SetUsername(user)
	}
	password := c.String("P")
	if password != "" {
		opts.SetPassword(password)
	}

	if host != "" {
		brokerUri := fmt.Sprintf("%s://%s:%d", scheme, host, port)
		log.Infof("Broker URI: %s", brokerUri)

		opts.AddBroker(brokerUri)
	}

	opts.SetAutoReconnect(true)
	return opts, nil
}
