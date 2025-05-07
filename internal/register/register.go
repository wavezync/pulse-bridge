package register

import (
	"time"
	"wavezync/pulse-bridge/internal/cache"
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

	resultChan := make(chan ResultChanStruct)

	go func() {
		var err error
		startTime := time.Now()

		switch mntr.Type {
		case "http":
			err = monitor.HttpMonitor(mntr)
		case "database":
			err = monitor.DatabaseMonitor(mntr)
		}

		duration := time.Since(startTime)

		// Consecutive successes logic
		oldResponse, isExisting := cache.DefaultMonitorCache.GetMonitorStatus(mntr.Name)
		consecutiveSuccesses := 0
		if err == nil {
			if isExisting {
				consecutiveSuccesses = oldResponse.Metrics.ConsecutiveSuccesses + 1
			} else {
				consecutiveSuccesses = 1
			}
		} else {
			consecutiveSuccesses = 0
		}

		resultChan <- ResultChanStruct{
			err:                  err,
			mntr:                 mntr,
			duration:             duration,
			consecutiveSuccesses: consecutiveSuccesses,
		}
	}()

	select {
	case result := <-resultChan:
		createResponse(result)
	case <-timer.C:
		log.Warn().
			Str("monitor", mntr.Name).
			Msg("Operation timed out")
	}
}
