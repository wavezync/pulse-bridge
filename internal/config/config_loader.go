package config

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
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

	for _, monitor := range config.Monitors {
		jsonData, _ := json.MarshalIndent(monitor, "", "  ")
		log.Debug().
			Str("monitor", monitor.Name).
			RawJSON("config", jsonData).
			Msg("Loaded monitor configuration")
	}

	return &config, nil
}
