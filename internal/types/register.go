package types

import (
	"time"
	"wavezync/pulse-bridge/internal/config"
)

type ResultChanStruct struct {
	Err      *MonitorError
	Mntr     *config.Monitor
	Duration time.Duration
}
