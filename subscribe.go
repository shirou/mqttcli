package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func subscribe(c *cli.Context) {
	if c.Bool("d") {
		log.SetLevel(log.DebugLevel)
	}
	opts, err := NewOption(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if c.Bool("c") {
		opts.SetCleanSession(false)
	}

	client, err := connect(c, opts)
	if err != nil {
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

	err = client.Subscribe(topic, qos)
	if err != nil {
		log.Error(err)
	}
}
