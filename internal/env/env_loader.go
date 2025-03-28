package env

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func Get() *Config {
	return &AppConfig
}

func Init() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warn().Err(err).Msg("Error loading .env file, using environment variables")
	}

	loadConfig()

	log.Info().Msg("Environment configuration loaded successfully")

	return &AppConfig
}

func loadConfig() {
	AppConfig = Config{
		ConfigPath: GetEnv("ATHENA_CONFIG", ""),
		Host:       GetEnv("HOST", "0.0.0.0"),
		Port:       GetEnvInt("PORT", 8080),
	}
}
