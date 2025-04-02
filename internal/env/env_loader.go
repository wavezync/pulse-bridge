package env

import (
	"github.com/joho/godotenv"
)

func Get() *Config {
	return &AppConfig
}

func Init() *Config {
	godotenv.Load(".env")
	loadConfig()

	return &AppConfig
}

func loadConfig() {
	AppConfig = Config{
		ConfigPath: GetEnv("PULSE_BRIDGE_CONFIG", ""),
		Host:       GetEnv("HOST", "0.0.0.0"),
		Port:       GetEnvInt("PORT", 8080),
	}
}
