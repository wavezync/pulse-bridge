package register

import (
	"wavezync/pulse-bridge/internal/config"

	"github.com/rs/zerolog/log"
)

func SetRegister(cfg *config.Config) {
	for _, monitor := range cfg.Monitors {
		switch monitor.Type {
		case "http":
			// Register HTTP monitor
			// registerHttpMonitor(monitor)
		case "database":
			// Register Database monitor
			// registerDatabaseMonitor(monitor)
		default:
			log.Error().
				Str("type", monitor.Type).
				Str("monitor", monitor.Name).
				Msg("Unknown monitor type")
		}
	}
}
