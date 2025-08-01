package register

import (
	"time"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/monitor"
	"wavezync/pulse-bridge/internal/types"

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
func runMonitorWorker(mntr *config.Monitor, monitoringTimerFn ...func(*config.Monitor)) {
	timerFn := monitoringTimer
	if len(monitoringTimerFn) > 0 {
		timerFn = monitoringTimerFn[0]
	}
	duration, err := time.ParseDuration(mntr.Interval)
	if err != nil {
		log.Error().
			Err(err).
			Str("monitor", mntr.Name).
			Msg("Error parsing duration, monitor not started")
		return
	}

	// Run the monitor check immediately first
	timerFn(mntr)

	// Then set up the ticker for subsequent checks
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		timerFn(mntr)
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

	resultChan := make(chan types.ResultChanStruct)

	go func() {
		var mntrErr *types.MonitorError
		startTime := time.Now()

		switch mntr.Type {
		case "http":
			mntrErr = monitor.HttpMonitor(mntr)
		case "database":
			mntrErr = monitor.DatabaseMonitor(mntr)
		}

		duration := time.Since(startTime)

		resultChan <- types.ResultChanStruct{
			Err:      mntrErr,
			Mntr:     mntr,
			Duration: duration,
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
