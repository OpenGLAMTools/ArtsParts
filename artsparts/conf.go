package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	ServerPort    string `yaml:"server_port,omitempty"`
	SessionSecret string `yaml:"session_secret,omitempty"`
	LogLevel      uint32 `yaml:"log_level"`
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
