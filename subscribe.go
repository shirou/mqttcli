package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func subscribe(c *cli.Context) error {
	setDebugLevel(c)
	opts, err := NewOption(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	opts.SetKeepAlive(time.Second * 60)

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

	return nil
}
