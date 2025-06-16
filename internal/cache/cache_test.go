package cache

import (
	"testing"
	"wavezync/pulse-bridge/internal/types"
)

func TestMonitorCache_SetAndGetMonitorStatus(t *testing.T) {
	cache := NewMonitorCache()
	resp := types.MonitorResponse{Service: "test", Status: types.StatusHealthy}

	cache.SetMonitorStatus("test", resp)
	got, exists := cache.GetMonitorStatus("test")
	if !exists {
		t.Fatalf("expected monitor status to exist")
	}
	if got != resp {
		t.Errorf("got %+v, want %+v", got, resp)
	}
}

func TestMonitorCache_GetMonitorStatus_NotExists(t *testing.T) {
	cache := NewMonitorCache()
	_, exists := cache.GetMonitorStatus("notfound")
	if exists {
		t.Errorf("expected monitor status to not exist")
	}
}

func TestMonitorCache_GetAllMonitorStatus(t *testing.T) {
	cache := NewMonitorCache()
	cache.SetMonitorStatus("a", types.MonitorResponse{Service: "a", Status: types.StatusHealthy})
	cache.SetMonitorStatus("b", types.MonitorResponse{Service: "b", Status: types.StatusUnhealthy})

	all := cache.GetAllMonitorStatus()
	if len(all) != 2 {
		t.Errorf("expected 2 monitor statuses, got %d", len(all))
	}
	names := map[string]bool{}
	for _, resp := range all {
		names[resp.Service] = true
	}
	if !names["a"] || !names["b"] {
		t.Errorf("expected monitor services 'a' and 'b' in results")
	}
}
