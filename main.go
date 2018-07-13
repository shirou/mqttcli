package main

import (
	"bufio"
	stdlog "log"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	colorable "github.com/mattn/go-colorable"
	"github.com/urfave/cli"
)

var usage = `
Usage here
`

var version string

func init() {
	log.SetLevel(log.WarnLevel)
	log.SetOutput(colorable.NewColorableStdout())
}

// connects MQTT broker
func connect(c *cli.Context, opts *MQTT.ClientOptions, subscribed map[string]byte) (*MQTTClient, error) {
	willPayload := c.String("will-payload")
	willQoS := c.Int("will-qos")
	willRetain := c.Bool("will-retain")
	willTopic := c.String("will-topic")
	if willPayload != "" && willTopic != "" {
		opts.SetWill(willTopic, willPayload, byte(willQoS), willRetain)
	}

	client := &MQTTClient{Opts: opts}
	client.lock = new(sync.Mutex)
	client.Subscribed = subscribed

	opts.SetOnConnectHandler(client.SubscribeOnConnect)
	opts.SetConnectionLostHandler(client.ConnectionLost)

	_, err := client.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func pubsub(c *cli.Context) error {
	setDebugLevel(c)
	opts, err := NewOption(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	qos := c.Int("q")
	subtopic := c.String("sub")
	if subtopic == "" {
		log.Errorf("Please specify sub topic")
		os.Exit(1)
	}
	log.Infof("Sub Topic: %s", subtopic)
	pubtopic := c.String("pub")
	if pubtopic == "" {
		log.Errorf("Please specify pub topic")
		os.Exit(1)
	}
	log.Infof("Pub Topic: %s", pubtopic)
	retain := c.Bool("r")

	subscribed := map[string]byte{
		subtopic: byte(0),
	}

	client, err := connect(c, opts, subscribed)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	go func() {
		// Read from Stdin and publish
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err = client.Publish(pubtopic, []byte(scanner.Text()), qos, retain, false)
			if err != nil {
				log.Error(err)
			}
		}
	}()

	// while loop
	for {
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "mqttcli"
	app.Usage = usage
	app.Version = version

	commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			Usage:  "mqtt host to connect to. Defaults to localhost",
			EnvVar: "MQTT_HOST"},
		cli.IntFlag{
			Name:   "p, port",
			Value:  1883,
			Usage:  "network port to connect to. Defaults to 1883",
			EnvVar: "MQTT_PORT"},
		cli.StringFlag{
			Name:   "u,user",
			Value:  "",
			Usage:  "provide a username",
			EnvVar: "MQTT_USERNAME"},
		cli.StringFlag{
			Name:   "P,password",
			Value:  "",
			Usage:  "provide a password",
			EnvVar: "MQTT_PASSWORD"},
		cli.StringFlag{
			Name:  "t",
			Value: "",
			Usage: "mqtt topic to publish to.",
		},
		cli.IntFlag{
			Name:  "q",
			Value: 0,
			Usage: "QoS",
		},
		cli.StringFlag{
			Name:  "cafile",
			Value: "",
			Usage: "CA certificates",
		},
		cli.StringFlag{
			Name:  "cert",
			Value: "",
			Usage: "Client certificates",
		},
		cli.StringFlag{
			Name:  "key",
			Value: "",
			Usage: "Client private key",
		},
		cli.StringFlag{
			Name:  "i",
			Value: "",
			Usage: "ClientiId. Defaults random.",
		},
		cli.StringFlag{
			Name:  "m",
			Value: "test message",
			Usage: "Message body",
		},
		cli.BoolFlag{
			Name:  "r",
			Usage: "message should be retained.",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "enable debug messages",
		},
		cli.BoolFlag{
			Name:  "dd",
			Usage: "enable debug messages",
		},
		cli.BoolFlag{
			Name:  "ddd",
			Usage: "enable debug messages",
		},
		cli.BoolFlag{
			Name:  "dddd",
			Usage: "enable debug messages",
		},
		cli.BoolFlag{
			Name:  "ddddd",
			Usage: "enable debug messages",
		},
		cli.BoolFlag{
			Name:  "insecure",
			Usage: "do not check that the server certificate",
		},
		cli.StringFlag{
			Name:   "conf",
			Value:  "~/.mqttcli.cfg",
			Usage:  "config file path",
			EnvVar: "MQTTCLI_CONFPATH"},
		cli.StringFlag{
			Name:  "will-payload",
			Value: "",
			Usage: "payload for the client Will",
		},
		cli.IntFlag{
			Name:  "will-qos",
			Value: 0,
			Usage: "QoS level for the client Will",
		},
		cli.BoolFlag{
			Name:  "will-retain",
			Usage: "if given, make the client Will retained",
		},
		cli.StringFlag{
			Name:  "will-topic",
			Value: "",
			Usage: "the topic on which to publish the client Will",
		},
	}
	pubFlags := append(commonFlags,
		cli.BoolFlag{
			Name:  "s",
			Usage: "read message from stdin, sending line by line as a message",
		},
	)
	subFlags := append(commonFlags,
		cli.BoolFlag{
			Name:  "c",
			Usage: "disable 'clean session'",
		},
	)
	pubsubFlags := append(commonFlags,
		cli.StringFlag{
			Name:  "pub",
			Usage: "publish topic",
		},
		cli.StringFlag{
			Name:  "sub",
			Usage: "subscribe topic",
		},
	)

	app.Commands = []cli.Command{
		{
			Name:   "pub",
			Usage:  "publish",
			Flags:  pubFlags,
			Action: publish,
		},
		{
			Name:   "sub",
			Usage:  "subscribe",
			Flags:  subFlags,
			Action: subscribe,
		},
		{
			Name:   "pubsub",
			Usage:  "subscribe and publish",
			Flags:  pubsubFlags,
			Action: pubsub,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}

func setDebugLevel(c *cli.Context) {
	if c.Bool("d") {
		log.SetLevel(log.DebugLevel)
	} else if c.Bool("dd") {
		log.SetLevel(log.DebugLevel)
		MQTT.WARN = stdlog.New(os.Stdout, "", stdlog.LstdFlags)
	} else if c.Bool("ddd") {
		log.SetLevel(log.DebugLevel)
		MQTT.DEBUG = stdlog.New(os.Stdout, "", stdlog.LstdFlags)
	}
}
