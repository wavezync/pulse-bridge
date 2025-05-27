package register

import (
	"time"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/types"
)

type ResultChanStruct struct {
	err                  *types.MonitorError
	mntr                 *config.Monitor
	duration             time.Duration
}
