package main

import (
	"encoding/json"
	"testing"
)

func Test_ConfigUnmarshalJSON(t *testing.T) {
	c := Config{}

	j := `{"host": "h", "port": 1883}`
	err := json.Unmarshal([]byte(j), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Port != 1883 {
		t.Error("port(int) is not correctly read")
	}
	if c.Host != "h" || c.UserName != "" || c.Password != "" {
		t.Error("parse failed")
	}

	j = `{"host": "h", "port": "1883", "username": "u"}`
	err = json.Unmarshal([]byte(j), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Port != 1883 {
		t.Error("port(string) is not correctly read")
	}
	if c.Host != "h" || c.UserName != "u" || c.Password != "" {
		t.Error("parse failed with username set")
	}
}

func Test_ConfigCert(t *testing.T) {
	c := Config{}

	j := `{"caCert": "ca", "clientCert": "client", "privateKey": "key"}`
	err := json.Unmarshal([]byte(j), &c)
	if err != nil {
		t.Error(err)
	}
	if c.CaCert != "ca" || c.ClientCert != "client" || c.PrivateKey != "key" {
		t.Errorf("parse failed, %+v", c)
	}
}
