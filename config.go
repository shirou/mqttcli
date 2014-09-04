package main

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigJson struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func readFromConfigFile(path string) (ConfigJson, error) {
	ret := ConfigJson{}

	b, err := ioutil.ReadFile("/tmp/dat")
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
