package handler

import (
	"encoding/json"
	"net/http"

	"strings"
	"wavezync/pulse-bridge/internal/cache"

	"github.com/rs/zerolog/log"
)

func MonitorServices(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("remote_addr", r.RemoteAddr).
		Str("method", r.Method).
		Msg("Monitor services status requested")

	w.Header().Set("Content-Type", "application/json")
	statuses := cache.DefaultMonitorCache.GetAllMonitorStatus()
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
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
