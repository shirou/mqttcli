package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var log = logrus.New()
var usage = `
Usage here
`

func initFunc() {
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.Debug
}

// connects MQTT broker
func connect(c *cli.Context) (*MQTTClient, error){
	host := c.String("host")
	port := c.Int("p")
	clientId := c.String("i")

	cafile := c.String("cafile")
	scheme := "tcp"
	if cafile != "" {
		scheme = "ssl"
	}

	user := c.String("u")
	password := c.String("P")

	brokerUri := fmt.Sprintf("%s://%s:%d", scheme, host, port)
	log.Info("Broker URI: %s", brokerUri)

	client := NewMQTTClient()

	log.Debug("Connecting...")
	_, err := client.Connect(brokerUri, clientId, user, password)
	if err != nil{
		return nil, err
	}
	log.Debug("Connected")

	return client, nil
}


func publish(c *cli.Context) {
	client, err := connect(c)
	if err != nil{
		log.Error(err)
		os.Exit(1)
	}

	qos := c.Int("q")
	topic := c.String("t")
	if topic == "" {
		log.Errorf("Please specify topic")
		os.Exit(1)
	}
	log.Infof("Topic: %s", topic)

	payload := c.String("m")
	retain := c.Bool("r")
	log.Infof("Retain: %t", retain)

	client.Publish(topic, []byte(payload), qos, retain)

	log.Debug("Published")
}

func main() {
	initFunc()

	app := cli.NewApp()
	app.Name = "mqttcli"
	app.Usage = usage
	app.Commands = []cli.Command{
		{
			Name:  "publish",
			Usage: "publish",
			Flags: []cli.Flag{
				cli.StringFlag{"host", "test.mosquitto.org", "Broker IP address"},
				cli.IntFlag{"p", 1883, "Broker Port"},
				cli.StringFlag{"t", "", "Topic"},
				cli.IntFlag{"q", 0, "QoS"},
				cli.StringFlag{"cafile", "", "CA file"},
				cli.StringFlag{"u", "", "username"},
				cli.StringFlag{"P", "", "password"},
				cli.StringFlag{"i", "client0", "ClientiId"},
				cli.StringFlag{"m", "test message", "Message body"},
				cli.BoolFlag{"r", "Retain flag"},
			},
			Action: publish,
		},
	}
	app.Run(os.Args)
}
