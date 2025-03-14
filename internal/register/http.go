package register

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"wavezync/pulse-bridge/internal/config"
)

func HttpMonitor(monitor *config.Monitor) (string, error) {
	timeout, err := time.ParseDuration(monitor.Timeout)
	if err != nil {
		return "", fmt.Errorf("invalid timeout format: %w", err)
	}

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(strings.ToUpper(monitor.Http.Method), monitor.Http.Url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range monitor.Http.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var jsonResp map[string]interface{}
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	status, ok := jsonResp["status"].(string)

	if status != "ok" || !ok {
		return "", fmt.Errorf("service reported non-ok status: %s", status)
	}

	return "Monitor check successful", nil
}
