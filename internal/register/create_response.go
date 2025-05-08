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
		log.Info().
			Str("monitor", result.mntr.Name).
			Msg("Monitor check failed")
	} else {
		log.Info().
			Str("monitor", result.mntr.Name).
			Msg("Monitor check successful")
	}

	oldResponse, isExisting := cache.DefaultMonitorCache.GetMonitorStatus(result.mntr.Name)

	lastSuccess := getLastSuccess(result.err, isExisting, oldResponse, time.Now().String())
	consecutiveSuccesses := getConsecutiveSuccesses(result.err, isExisting, oldResponse)
	lastError := getLastError(result.err, isExisting, oldResponse)
	status := statusFromError(result.err)

	newResponse = types.MonitorResponse{
		Service:     result.mntr.Name,
		Status:      status,
		Type:        result.mntr.Type,
		LastCheck:   time.Now().String(),
		LastSuccess: lastSuccess,
		Metrics: types.Metrics{
			ResponseTimeMs:       int(result.duration.Milliseconds()),
			CheckInterval:        result.mntr.Interval,
			ConsecutiveSuccesses: consecutiveSuccesses,
		},
		LastError: lastError,
	}

	cache.DefaultMonitorCache.SetMonitorStatus(result.mntr.Name, newResponse)
	return newResponse
}

func getLastSuccess(err *types.MonitorError, isExisting bool, oldResponse types.MonitorResponse, currentTime string) string {
	if err == nil {
		return currentTime
	} else if isExisting && oldResponse.Status == types.StatusHealthy {
		return oldResponse.LastCheck
	} else if isExisting {
		return oldResponse.LastSuccess
	}
	return ""
}

func getConsecutiveSuccesses(err *types.MonitorError, isExisting bool, oldResponse types.MonitorResponse) int {
	if err == nil {
		if isExisting {
			return oldResponse.Metrics.ConsecutiveSuccesses + 1
		}
		return 1
	}
	return 0
}

func getLastError(err *types.MonitorError, isExisting bool, oldResponse types.MonitorResponse) string {
	if err != nil {
		return err.Error()
	} else if isExisting && oldResponse.LastError != "" {
		return oldResponse.LastError
	}
	return ""

}

func statusFromError(err *types.MonitorError) types.Status {
	if err == nil {
		return types.StatusHealthy
	} else if types.IsClientError(err) {
		return types.StatusUnhealthy
	} else if types.IsConfigError(err) {
		return types.StatusUnknown
	}
	return types.StatusUnknown
}
