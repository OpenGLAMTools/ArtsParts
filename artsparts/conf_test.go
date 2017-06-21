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
		"path/to/src/root",
		5,
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
