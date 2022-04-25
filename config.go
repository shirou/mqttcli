package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	simpleJson "github.com/bitly/go-simplejson"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const DefaultConfigFilePath = "~/.mqttcli.cfg"

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`

	CaCert     string `json:"caCert"`
	ClientCert string `json:"clientCert"`
	PrivateKey string `json:"privateKey"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	js, err := simpleJson.NewJson(data)
	if err != nil {
		return err
	}
	if c.Host, err = js.Get("host").String(); err != nil {
		c.Host = ""
	}
	// Port can be string either int
	if c.Port, err = js.Get("port").Int(); err != nil {
		p, err := js.Get("port").String()
		c.Port, err = strconv.Atoi(p)
		if err != nil {
			c.Port = 0
		}
	}
	if c.UserName, err = js.Get("username").String(); err != nil {
		c.UserName = ""
	}
	if c.Password, err = js.Get("password").String(); err != nil {
		c.Password = ""
	}
	if c.CaCert, err = js.Get("caCert").String(); err != nil {
		c.CaCert = ""
	}
	if c.ClientCert, err = js.Get("clientCert").String(); err != nil {
		c.ClientCert = ""
	}
	if c.PrivateKey, err = js.Get("privateKey").String(); err != nil {
		c.PrivateKey = ""
	}
	return nil
}

func readFromConfigFile(path string) (Config, error) {
	ret := Config{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func existsDefaultConfigFile() (string, bool) {
	p := strings.Replace(DefaultConfigFilePath, "~", UserHomeDir(), 1)
	if _, err := os.Stat(p); err == nil {
		return p, true
	} else {
		return p, false
	}
}

func getSettingsFromFile(confPath string, opts *MQTT.ClientOptions) error {
	ret, err := readFromConfigFile(confPath)
	if err != nil {
		log.Error(err)
		return err
	}

	tlsConfig, ok, err := makeTlsConfig(ret.CaCert, ret.ClientCert, ret.PrivateKey, false)
	if err != nil {
		return err
	}
	if ok {
		opts.SetTLSConfig(tlsConfig)
	}

	if ret.Host != "" {
		if ret.Port == 0 {
			ret.Port = 1883
		}
		scheme := "tcp"
		if ret.Port == 8883 {
			scheme = "ssl"
		}
		brokerUri := fmt.Sprintf("%s://%s:%d", scheme, ret.Host, ret.Port)
		log.Infof("Broker URI: %s", brokerUri)
		opts.AddBroker(brokerUri)
	}

	if ret.UserName != "" {
		opts.SetUsername(ret.UserName)
	}
	if ret.Password != "" {
		opts.SetPassword(ret.Password)
	}
	return nil
}
