package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseConf(t *testing.T) {
	b, err := ioutil.ReadFile("default.conf.yml")
	if err != nil {
		t.Error("Error reading default conf file", err)
	}
	conf, err := parseConf(b)
	exp := Conf{
		":3000",
		"secret_and_long_string",
		"../test",
		"http://test.com",
		5,
		map[string]string{
			"TWITTER_KEY":         "KEY123",
			"TWITTER_SECRET":      "SECRET123",
			"ACCESS_TOKEN":        "ACCESS",
			"ACCESS_TOKEN_SECRET": "ACCESS_SECRET",
		},
		map[string]string{
			"Impressum":   "This is the impressum",
			"AnotherPage": "This is another text",
		},
	}
	if !reflect.DeepEqual(conf, exp) {
		t.Error("Conf parsing is not correct", exp, conf)
	}
}

func TestLoadConfError(t *testing.T) {
	_, err := loadConf("notexist")
	if err == nil {
		t.Error("Expect an error")
	}
}

func TestEnsureConf(t *testing.T) {

}
