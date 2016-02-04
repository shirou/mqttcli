package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func subscribe(c *cli.Context) {
	setDebugLevel(c)
	opts, err := NewOption(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Setting KeepAlive to 0 disables it. Paho MQTT client is currently broken
	// and does not send ping when subscribing only.
	// TODO set KeepAlive to a real value (60s?) when this change is merged:
	// https://git.eclipse.org/r/#/c/65850/
	opts.SetKeepAlive(time.Duration(0))

	if c.Bool("c") {
		clientId := c.String("i")
		if clientId == "" {
			log.Warn("clean Flag does not work without client id")
		}

		opts.SetCleanSession(false)
	}

	qos := c.Int("q")
	topic := c.String("t")
	if topic == "" {
		log.Errorf("Please specify topic")
		os.Exit(1)
	}
	log.Infof("Topic: %s", topic)

	subscribed := map[string]byte{
		topic: byte(qos),
	}

	_, err = connect(c, opts, subscribed)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// loops forever
	for {
		time.Sleep(time.Second)
	}

}
