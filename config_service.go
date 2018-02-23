package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func (c *ConfigurationFile) GetConfig() *ConfigurationFile {
	configRaw, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configRaw, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}