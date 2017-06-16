package main

import "io/ioutil"
import "gopkg.in/yaml.v2"

type Conf struct {
	ServerPort    string
	TwitterKey    string
	TwitterSecret string
}

func loadConf() Conf {
	b, _ := ioutil.ReadFile(".conf.yml")
	c := Conf{}
	yaml.Unmarshal(b, c)
	return c
}
