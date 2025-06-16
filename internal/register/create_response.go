package register

import (
	"time"
	"wavezync/pulse-bridge/internal/cache"
	"wavezync/pulse-bridge/internal/types"

	"github.com/rs/zerolog/log"
)

func createResponse(result types.ResultChanStruct) types.MonitorResponse {
	var newResponse types.MonitorResponse

	if result.Err != nil {
		log.Info().
			Str("monitor", result.Mntr.Name).
			Msg("Monitor check failed")
	} else {
		log.Info().
			Str("monitor", result.Mntr.Name).
			Msg("Monitor check successful")
	}

	oldResponse, isExisting := cache.DefaultMonitorCache.GetMonitorStatus(result.Mntr.Name)

	lastSuccess := getLastSuccess(result.Err, isExisting, oldResponse, time.Now().String())
	consecutiveSuccesses := getConsecutiveSuccesses(result.Err, isExisting, oldResponse)
	lastError := getLastError(result.Err, isExisting, oldResponse)
	status := statusFromError(result.Err)

	newResponse = types.MonitorResponse{
		Service:     result.Mntr.Name,
		Status:      status,
		Type:        result.Mntr.Type,
		LastCheck:   time.Now().String(),
		LastSuccess: lastSuccess,
		Metrics: types.Metrics{
			ResponseTimeMs:       int(result.Duration.Milliseconds()),
			CheckInterval:        result.Mntr.Interval,
			ConsecutiveSuccesses: consecutiveSuccesses,
		},
		LastError: lastError,
	}

	cache.DefaultMonitorCache.SetMonitorStatus(result.Mntr.Name, newResponse)
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
