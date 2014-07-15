package main

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"crypto/tls"
)

type MQTTClient struct {
	Opts   *MQTT.ClientOptions
	Client *MQTT.MqttClient
}

func NewMQTTClient() *MQTTClient {
	return &MQTTClient{}
}

// connect to MQTT broker
func (m *MQTTClient) Connect(brokerUri string, clientId string,
	user string, password string) (*MQTT.MqttClient, error) {

	m.Opts = MQTT.NewClientOptions()

	m.Opts.SetBroker(brokerUri)
	m.Opts.SetClientId(clientId)
	m.Opts.SetTraceLevel(MQTT.Critical)
//	m.Opts.SetTraceLevel(MQTT.Verbose)
	if user != "" {
		m.Opts.SetUsername(user)
	}
	if password != "" {
		m.Opts.SetPassword(password)
	}

	insecure := true
	if insecure {
		tlsConfig := &tls.Config{InsecureSkipVerify: true,}
		m.Opts.SetTlsConfig(tlsConfig)
	}

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

	//      receipt := m.Client.PublishMessage(msg.Destination, mqttmsg)
	receipt := m.Client.PublishMessage(topic, mqttmsg)
	<-receipt

	return nil
}
