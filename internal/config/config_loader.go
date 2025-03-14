package config

import (
	"os"
	"wavezync/pulse-bridge/internal/env"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var config Config

func Get() *Config {
	return &config
}

func Init(configPath string, envConfig *env.Config) (*Config, error) {

	if configPath == "" {
		if envConfig.ConfigPath != "" {
			configPath = envConfig.ConfigPath
		} else {
			configPath = "config.yml"
		}
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", configPath).Msg("Failed to read config file")
	}

	err = yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		return nil, err
	}

	err = ValidateConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
