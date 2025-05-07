package cache

import (
	"sync"
	"wavezync/pulse-bridge/internal/types"
)

type MonitorCache struct {
	cache map[string]types.MonitorResponse
	lock  sync.RWMutex
}

func NewMonitorCache() *MonitorCache {
	return &MonitorCache{
		cache: make(map[string]types.MonitorResponse),
	}
}

func (mc *MonitorCache) SetMonitorStatus(name string, response types.MonitorResponse) {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.cache[name] = response
}

func (mc *MonitorCache) GetMonitorStatus(name string) (types.MonitorResponse, bool) {
	mc.lock.RLock()
	defer mc.lock.RUnlock()
	resp, exists := mc.cache[name]
	return resp, exists
}

func (mc *MonitorCache) GetAllMonitorStatus() []types.MonitorResponse {
	mc.lock.RLock()
	defer mc.lock.RUnlock()

	statuses := make([]types.MonitorResponse, 0, len(mc.cache))
	for _, status := range mc.cache {
		statuses = append(statuses, status)
	}

	return statuses
}

var DefaultMonitorCache = NewMonitorCache()
