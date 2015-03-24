package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var MaxClientIdLen = 8

type MQTTClient struct {
	Client *MQTT.Client
	Opts   *MQTT.ClientOptions
}

// Connects connect to the MQTT broker with Options.
func (m *MQTTClient) Connect() (*MQTT.Client, error) {
	m.Client = MQTT.NewClient(m.Opts)
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

func onMessageReceived(client *MQTT.Client, message MQTT.Message) {
	log.Infof("topic:%s  / msg:%s", message.Topic(), message.Payload())
	fmt.Println(string(message.Payload()))
}

func (m *MQTTClient) Subscribe(topic string, qos int) error {
	token := m.Client.Subscribe(topic, byte(qos), onMessageReceived)
	if token.Error() != nil {
		return token.Error()
	}

	for {
		time.Sleep(1 * time.Second)
	}
	return nil
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
	return opts, nil
}
