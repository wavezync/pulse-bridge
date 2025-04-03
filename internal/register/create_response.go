package register

import (
	"wavezync/pulse-bridge/internal/cache"

	"github.com/rs/zerolog/log"
)

func createResponse(resultChan ResultChanStruct) cache.MonitorResponse {
	status := "healthy"
	var errorMsg string

	if resultChan.err != nil {
		log.Error().
			Err(resultChan.err).
			Str("monitor", resultChan.mntr.Name).
			Msg("Monitor check failed")
		status = "unhealthy"
		errorMsg = resultChan.err.Error()
	} else {
		log.Info().
			Str("monitor", resultChan.mntr.Name).
			Msg("Monitor check successful")
	}

	response := cache.MonitorResponse{
		Service:     resultChan.mntr.Name,
		Status:      status,
		Type:        resultChan.mntr.Type,
		LastCheck:   resultChan.lastCheck.String(),
		LastSuccess: resultChan.lastCheck.String(),
		LastError:   errorMsg,
		Metrics: cache.Metrics{
			ResponseTimeMs:       int(resultChan.responseTime.Milliseconds()),
			CheckInterval:        resultChan.mntr.Interval,
			ConsecutiveSuccesses: 0,
		},
	}

	cache.SetMonitorStatus(resultChan.mntr.Name, response)

	return response
}
