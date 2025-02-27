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

	jcart, err := json.MarshalIndent(config, "", "  ")

	fmt.Println(string(jcart))

	return &config, nil
}
