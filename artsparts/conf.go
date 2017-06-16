package main

import "io/ioutil"
import "gopkg.in/yaml.v2"

type Conf struct {
	ServerPort    string `yaml:"server_port,omitempty"`
	TwitterKey    string `yaml:"twitter_key,omitempty"`
	TwitterSecret string `yaml:"twitter_secret,omitempty"`
}

func loadConf() (Conf, error) {
	b, _ := ioutil.ReadFile(".conf.yml")
	return parseConf(b)
}
func parseConf(b []byte) (Conf, error) {
	c := Conf{}
	err := yaml.Unmarshal(b, &c)
	return c, err
}
