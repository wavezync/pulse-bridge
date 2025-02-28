package config

import (
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

	err = ValidateConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
