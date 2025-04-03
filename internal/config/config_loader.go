package config

import (
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

	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Info().
		Interface("config", config).
		Msg("Configuration loaded successfully")

	return &config, nil
}
