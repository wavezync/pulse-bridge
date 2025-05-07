package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"wavezync/pulse-bridge/internal/cache"
	"wavezync/pulse-bridge/internal/config"

	"github.com/rs/zerolog/log"
)

func MonitorServices(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("remote_addr", r.RemoteAddr).
		Str("method", r.Method).
		Msg("Monitor services status requested")

	w.Header().Set("Content-Type", "application/json")

	allStatuses := cache.DefaultMonitorCache.GetAllMonitorStatus()

	statusMap := make(map[string]int)
	for i, status := range allStatuses {
		statusMap[status.Service] = i
	}

	cfg := config.Get()

	orderedStatuses := make([]interface{}, 0, len(allStatuses))
	for _, monitor := range cfg.Monitors {
		if idx, exists := statusMap[monitor.Name]; exists {
			orderedStatuses = append(orderedStatuses, allStatuses[idx])
		}
	}

	for _, status := range allStatuses {
		found := false
		for _, monitor := range cfg.Monitors {
			if status.Service == monitor.Name {
				found = true
				break
			}
		}
		if !found {
			orderedStatuses = append(orderedStatuses, status)
		}
	}

	if err := json.NewEncoder(w).Encode(orderedStatuses); err != nil {
		log.Error().Err(err).Msg("Failed to encode monitor services response")
	}
}

func MonitorServiceByName(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("remote_addr", r.RemoteAddr).
		Str("method", r.Method).
		Msg("Monitor service by name requested")

	name := strings.TrimPrefix(r.URL.Path, "/monitor/services/")
	if name == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}

	resp, exists := cache.DefaultMonitorCache.GetMonitorStatus(name)
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("Failed to encode monitor service response")
	}
}
