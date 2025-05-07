package types

type MonitorResponse struct {
	Service     string     `json:"service"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
	LastCheck   string     `json:"last_check"`
	LastSuccess string     `json:"last_success"`
	Metrics     Metrics    `json:"metrics"`
	LastError   string    `json:"last_error"`
}

type Metrics struct {
	ResponseTimeMs       int    `json:"response_time_ms"`
	CheckInterval       string `json:"check_interval"`
	ConsecutiveSuccesses int    `json:"consecutive_successes"`
}
