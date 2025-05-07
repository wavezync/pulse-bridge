package register

import (
	"time"
	"wavezync/pulse-bridge/internal/cache"
	"wavezync/pulse-bridge/internal/types"

	"github.com/rs/zerolog/log"
)

func createResponse(result ResultChanStruct) types.MonitorResponse {
	var newResponse types.MonitorResponse

	if result.err != nil {
		log.Error().
			Err(result.err).
			Str("monitor", result.mntr.Name).
			Msg("Monitor check failed")
	} else {
		log.Info().
			Str("monitor", result.mntr.Name).
			Msg("Monitor check successful")
	}

	oldResponse, isExisting := cache.DefaultMonitorCache.GetMonitorStatus(result.mntr.Name)

	var lastSuccess string
	if result.err == nil {
		lastSuccess = time.Now().String()
	} else if isExisting && oldResponse.Status == "healthy" {
		lastSuccess = oldResponse.LastCheck
	} else if isExisting {
		lastSuccess = oldResponse.LastSuccess
	} else {
		lastSuccess = ""
	}

	// Consecutive successes logic
	consecutiveSuccesses := 0
	if result.err == nil {
		if isExisting {
			consecutiveSuccesses = oldResponse.Metrics.ConsecutiveSuccesses + 1
		} else {
			consecutiveSuccesses = 1
		}
	} else {
		consecutiveSuccesses = 0
	}

	newResponse = types.MonitorResponse{
		Service:     result.mntr.Name,
		Status:      statusFromError(result.err),
		Type:        result.mntr.Type,
		LastCheck:   time.Now().String(),
		LastSuccess: lastSuccess,
		Metrics: types.Metrics{
			ResponseTimeMs:       int(result.duration.Milliseconds()),
			CheckInterval:        result.mntr.Interval,
			ConsecutiveSuccesses: consecutiveSuccesses,
		},
		LastError: errorString(result.err),
	}

	cache.DefaultMonitorCache.SetMonitorStatus(result.mntr.Name, newResponse)
	return newResponse
}

func statusFromError(err error) string {
	if err == nil {
		return "healthy"
	}
	return "unhealthy"
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
