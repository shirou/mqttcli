package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/Sirupsen/logrus"
)

const DefaultConfigFile = ".mqttcli.cfg" // Under HOME

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
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
func getSettingsFromFile(p string, opts *MQTT.ClientOptions) error {
	confPath := ""

	if p == "~/.mqtt.cfg" || p == "" {
		home := UserHomeDir()
		confPath = path.Join(home, DefaultConfigFile)

		_, err := os.Stat(confPath)
		if os.IsNotExist(err) {
			return err
		}
	} else {
		confPath = p
	}

	ret, err := readFromConfigFile(confPath)
	if err != nil {
		return err
	}
	if ret.Host != "" {
		if ret.Port == "" {
			ret.Port = "1883"
		}
		scheme := "tcp" // FIXME:
		brokerUri := fmt.Sprintf("%s://%s:%s", scheme, ret.Host, ret.Port)
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
