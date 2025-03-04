package env

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func Get() *Config {
	return &AppConfig
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Error loading .env file, using environment variables")
	}

	loadConfig()

	log.Info().Msg("Environment configuration loaded successfully")

	return &AppConfig
}

func loadConfig() {
	AppConfig = Config{
		ConfigPath: GetEnv("PULSE_BRIDGE_CONFIG", ""),
		Host:       GetEnv("PULSE_BRIDGE_HOST", "0.0.0.0"),
		Port:       GetEnvInt("PULSE_BRIDGE_PORT", 8080),
	}
}
