package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var conf Conf

func init() {
	var err error
	conf, err = loadConf(confFile)
	if err != nil {
		log.Fatal("Error loading conf", err)
	}
}

type Conf struct {
	ServerPort    string `yaml:"server_port,omitempty"`
	SessionSecret string `yaml:"session_secret,omitempty"`
}

func loadConf(f string) (Conf, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return Conf{}, err
	}
	return parseConf(b)
}
func parseConf(b []byte) (Conf, error) {
	c := Conf{}
	err := yaml.Unmarshal(b, &c)
	return c, err
}
