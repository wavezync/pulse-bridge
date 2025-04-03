package pulsebridge

import (
	"fmt"
	"net/http"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/env"
	"wavezync/pulse-bridge/internal/handler"
	"wavezync/pulse-bridge/internal/register"

	"github.com/rs/zerolog/log"
)

func Run(envConfig *env.Config) error {
	cfg, err := config.Init(envConfig.ConfigPath, envConfig)

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize configuration")
		return err
	}

	register.SetRegister(cfg)

	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/monitor/services", handler.Status)
	http.HandleFunc("/monitor/services/{serviceName}", handler.ServiceStatus)

	serverAddr := fmt.Sprintf("%s:%d", envConfig.Host, envConfig.Port)
	log.Info().Str("address", serverAddr).Msg("Starting server")

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Error().Err(err).Msg("Failed to start server")
		return err
	}

	return nil
}
