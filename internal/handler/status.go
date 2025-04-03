package handler

import (
	"encoding/json"
	"net/http"
	"wavezync/pulse-bridge/internal/cache"

	"github.com/rs/zerolog/log"
)

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statuses := cache.GetAllMonitorStatus()

	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ServiceStatus retrieves the status of a specific service
func ServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get service name from URL path parameter
	serviceName := r.PathValue("serviceName")

	if serviceName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Service name is required"})
		return
	}

	status, exists := cache.GetMonitorStatus(serviceName)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Service not found"})
		return
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
