package register

import (
	"time"
	"wavezync/pulse-bridge/internal/config"
)

type ResultChanStruct struct {
	err                  error
	mntr                 *config.Monitor
	duration             time.Duration
	consecutiveSuccesses int
}
