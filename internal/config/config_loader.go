package config

import (
	"encoding/json"
	"fmt"
	"wavezync/pulse-bridge/internal/env"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config Config

func Get() *Config {
	return &config
}

func Init(configPath string, envConfig *env.Config) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Str("path", configPath).Msg("Failed to read config file")
	}

	rawConfig := v.AllSettings()

	prettyJSON, err := json.MarshalIndent(rawConfig, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal config")
	} else {
		fmt.Println(string(prettyJSON))
	}

	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
