package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

var config Config

func Get() *Config {
	return &config
}

func Init(configFile string) (*Config, error) {

	err := yaml.Unmarshal([]byte(configFile), &config)

	if err != nil {
		return nil, err
	}

	for _, monitor := range config.Monitors {
		jcart, _ := json.MarshalIndent(monitor, "", "  ")

		fmt.Println(string(jcart))
	}

	return &config, nil
}
