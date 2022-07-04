package main

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// DefultConfigPath is default path for app config from rpm package
const DefultConfigPath = "/etc/test_quic_cfg.yaml"

type Config struct {
	Urls               []string `yaml:"urls"`
	MonFile            string   `yaml:"mon_file"`
	GoroutinesCount    int      `yaml:"goroutines"`
	ExpectedStatusCode string   `yaml:"expected_status_code"`
}

// GetConfig reading and parsing configuration yaml file
func (conf *Config) GetConfig(configPath string) {
	if configPath == "" {
		configPath = DefultConfigPath
	}
	yamlConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlConfig, &conf)
	if err != nil {
		fmt.Println("Unmarshal config error")
		log.Fatal(err)
	}
}
