package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Status(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("remote_addr", r.RemoteAddr).
		Str("method", r.Method).
		Msg("Status check requested")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
	}
}
