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
			Str("monitor", m.Name).
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
			Str("monitor", mntr.Name).
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
			Str("monitor", mntr.Name).
			Msg("Error parsing duration, monitor not started")
		return
	}

	timer := time.NewTimer(operationTimeout)
	defer timer.Stop()

	resultChan := make(chan struct {
		err error
	})

	go func() {
		var err error

		switch mntr.Type {
		case "http":
			err = monitor.HttpMonitor(mntr)
		case "database":
			err = monitor.DatabaseMonitor(mntr)
		}
		resultChan <- struct {
			err error
		}{err: err}
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
				Str("monitor", mntr.Name).
				Msg("Monitor check successful")
		}
	case <-timer.C:
		log.Warn().
			Str("monitor", mntr.Name).
			Msg("Operation timed out")
	}
}
