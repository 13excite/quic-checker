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
	Urls               []string `yaml:"urls"`
	GoroutinesCount    int      `yaml:"goroutines"`
	ExpectedStatusCode int      `yaml:"expected_status_code"`
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
	conf.ExpectedStatusCode = 200
	conf.Urls = []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.youtube.com",
		"https://www.google.com/1",
		"https://www.google.com/2",
		"https://www.google.com/3",
		"https://www.google.com/4",
		"https://www.google.com/5",
		"https://www.google.com/6",
		"https://www.google.com/7",
		"https://www.google.com/8",
		"https://www.google.com/9",
		"https://www.google.com/10",
		"https://www.google.com/11",
		"https://www.google.com/12",
		"https://www.google.com/13",
		"https://www.google.com/14",
		"https://www.google.com/15",
	}
}
