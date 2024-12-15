package config

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// DefultConfigPath is default path for app config from rpm package
const DefultConfigPath = "/etc/quic-checker.yaml"

// Config type represents a main configuration object
type Config struct {
	Urls            []URL `yaml:"urls"`
	GoroutinesCount int   `yaml:"goroutines"`
}

// Url type represents a URL object
type URL struct {
	URL              string
	ExpectStatusCode int
}

// type LoggerConfig struct {
// 	Color             bool     `yaml:"log_color"`
// 	DisableStacktrace bool     `yaml:"log_disable_stacktrace"`
// 	DevMode           bool     `yaml:"log_dev_mode"`
// 	DisableCaller     bool     `yaml:"log_disable_caller"`
// 	Level             string   `yaml:"log_level"`
// 	Encoding          string   `yaml:"log_encoding"`
// 	ErrorOutputPaths  []string `yaml:"log_err_output_paths"`
// }

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
	conf.GoroutinesCount = 3
	conf.Urls = []URL{
		{URL: "https://www.google.com", ExpectStatusCode: 200},
		{URL: "https://www.facebook.com", ExpectStatusCode: 200},
		{URL: "https://www.youtube.com", ExpectStatusCode: 200},
		{URL: "https://www.google.com/1", ExpectStatusCode: 200},
		{URL: "https://www.google.com/2", ExpectStatusCode: 200},
		{URL: "https://www.google.com/3", ExpectStatusCode: 200},
		{URL: "https://www.google.com/4", ExpectStatusCode: 200},
		{URL: "https://www.google.com/5", ExpectStatusCode: 200},
		{URL: "https://www.google.com/6", ExpectStatusCode: 200},
		{URL: "https://www.google.com/7", ExpectStatusCode: 200},
		{URL: "https://www.google.com/8", ExpectStatusCode: 200},
		{URL: "https://www.google.com/9", ExpectStatusCode: 200},
		{URL: "https://www.google.com/10", ExpectStatusCode: 200},
		{URL: "https://www.google.com/11", ExpectStatusCode: 200},
		{URL: "https://www.google.com/12", ExpectStatusCode: 200},
	}
}
