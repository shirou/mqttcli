package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var MaxClientIdLen = 8

type MQTTClient struct {
	Client *MQTT.MqttClient
	Opts   *MQTT.ClientOptions
}

// Connects connect to the MQTT broker with Options.
func (m *MQTTClient) Connect() (*MQTT.MqttClient, error) {
	m.Client = MQTT.NewClient(m.Opts)
	_, err := m.Client.Start()
	if err != nil {
		return nil, err
	}
	return m.Client, nil
}

func (m *MQTTClient) Publish(topic string, payload []byte, qos int, retain bool) error {
	mqttmsg := MQTT.NewMessage(payload)
	// FIXME: validate qos number
	mqttmsg.SetQoS(MQTT.QoS(qos))
	mqttmsg.SetRetainedFlag(retain)

	receipt := m.Client.PublishMessage(topic, mqttmsg)
	<-receipt

	return nil
}

func onMessageReceived(client *MQTT.MqttClient, message MQTT.Message) {
	log.Infof("topic:%s  / msg:%s", message.Topic(), message.Payload())
	fmt.Println(string(message.Payload()))
}

func (m *MQTTClient) Subscribe(topic string, qos int) error {
	topicFilter, err := MQTT.NewTopicFilter(topic, byte(qos))
	if err != nil {
		return err
	}
	_, err = m.Client.StartSubscription(onMessageReceived, topicFilter)
	if err != nil {
		return err
	}

	for {
		time.Sleep(1 * time.Second)
	}
	return nil
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
func NewOption(c *cli.Context) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions()

	host := c.String("host")
	port := c.Int("p")

	if host == "" && port == 0 {
		getSettingsFromFile(c.String("conf"), opts)
	}

	clientId := c.String("i")
	if clientId == "" {
		clientId = getRandomClientId()
	}
	opts.SetClientId(clientId)

	cafile := c.String("cafile")
	scheme := "tcp"
	if cafile != "" {
		scheme = "ssl"
	}
	insecure := true
	if insecure {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		opts.SetTlsConfig(tlsConfig)
	}

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
	return opts
}
