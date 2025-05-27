package monitor

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/types"
)

func HttpMonitor(monitor *config.Monitor) *types.MonitorError {
	timeout, err := time.ParseDuration(monitor.Timeout)
	if err != nil {
		return types.NewConfigError(fmt.Errorf("invalid timeout format: %w", err))
	}

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(strings.ToUpper(monitor.Http.Method), monitor.Http.Url, nil)
	if err != nil {
		return types.NewConfigError(fmt.Errorf("failed to create request: %w", err))
	}

	for key, value := range monitor.Http.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return types.NewClientError(fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NewClientError(fmt.Errorf("failed to read response: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		return types.NewClientError(fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}
