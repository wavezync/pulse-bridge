package types

type MonitorResponse struct {
	Service     string  `json:"service"`
	Status      Status  `json:"status"`
	Type        string  `json:"type"`
	LastCheck   string  `json:"last_check"`
	LastSuccess string  `json:"last_success"`
	Metrics     Metrics `json:"metrics"`
	LastError   string  `json:"last_error"`
}

type Metrics struct {
	ResponseTimeMs       int    `json:"response_time_ms"`
	CheckInterval        string `json:"check_interval"`
	ConsecutiveSuccesses int    `json:"consecutive_successes"`
}

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
	StatusUnknown   Status = "unknown"
)

func (s Status) String() string {
	return string(s)
}
