package config

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// DefultConfigPath is default path for app config from rpm package
const DefultConfigPath = "/etc/quic-checker.yaml"

type Config struct {
	Urls []string `yaml:"urls"`
	// MonFile            string   `yaml:"mon_file"`
	GoroutinesCount    int `yaml:"goroutines"`
	ExpectedStatusCode int `yaml:"expected_status_code"`
}

// GetConfig reading and parsing configuration yaml file
func (conf *Config) GetConfig(configPath string) {
	if configPath == "" {
		configPath = DefultConfigPath
	}
	yamlConfig, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlConfig, &conf)
	if err != nil {
		fmt.Println("Unmarshal config error")
		log.Fatal(err)
	}
}

func (conf *Config) Defaults() {
	conf.GoroutinesCount = 1
	conf.ExpectedStatusCode = 200
	conf.Urls = []string{"https://www.google.com"}
}
