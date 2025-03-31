package register

import (
	"time"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/monitor"

	"github.com/rs/zerolog/log"
)

func SetRegister(cfg *config.Config) {
	for _, monitor := range cfg.Monitors {
		m := monitor

		go runMonitorWorker(&m)

		log.Info().
			Str("name", m.Name).
			Str("interval", m.Interval).
			Msg("Monitor Initialized")
	}
}

// runMonitorWorker spawns a persistent worker that runs monitoring at specified intervals
func runMonitorWorker(mntr *config.Monitor) {
	duration, err := time.ParseDuration(mntr.Interval)
	if err != nil {
		log.Error().
			Err(err).
			Str("type", mntr.Type).
			Msg("Error parsing duration, monitor not started")
		return
	}

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		monitoringTimer(mntr)
	}
}

func monitoringTimer(mntr *config.Monitor) {
	operationTimeout, err := time.ParseDuration(mntr.Timeout)

	if err != nil {
		log.Error().
			Err(err).
			Str("name", mntr.Name).
			Msg("Error parsing duration, monitor not started")
		return
	}

	timer := time.NewTimer(operationTimeout)
	defer timer.Stop()

	resultChan := make(chan struct {
		message string
		err     error
	})

	go func() {
		var result string
		var err error

		switch mntr.Type {
		case "http":
			result, err = monitor.HttpMonitor(mntr)
		case "database":
			result, err = monitor.DatabaseMonitor(mntr)
		}
		resultChan <- struct {
			message string
			err     error
		}{message: result, err: err}
	}()

	select {
	case result := <-resultChan:
		if result.err != nil {
			log.Error().
				Err(result.err).
				Str("monitor", mntr.Name).
				Msg("Monitor check failed")
		} else {
			log.Info().
				Str("status", "ok").
				Str("name", mntr.Name)
		}
	case <-timer.C:
		log.Warn().
			Str("name", mntr.Name).
			Msg("Operation timed out")
	}
}
