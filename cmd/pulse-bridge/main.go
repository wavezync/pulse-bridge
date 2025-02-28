package main

import (
	"net/http"
	"os"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/handler"
	"wavezync/pulse-bridge/internal/register"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := "config.yml"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", configPath).Msg("Failed to read config file")
	}

	cfg, err := config.Init(string(configData))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize config")
	}

	register.SetRegister(cfg)

	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/status", handler.Status)

	log.Info().Msg("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
