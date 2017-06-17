package main

import "io/ioutil"
import "gopkg.in/yaml.v2"

type Conf struct {
	ServerPort    string `yaml:"server_port,omitempty"`
	SessionSecret string `yaml:"session_secret,omitempty"`
	TwitterKey    string `yaml:"twitter_key,omitempty"`
	TwitterSecret string `yaml:"twitter_secret,omitempty"`
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
