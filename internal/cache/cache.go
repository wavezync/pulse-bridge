package cache

import (
	"sync"
)

var (
	monitorCache     = make(map[string]MonitorResponse)
	monitorCacheLock = sync.RWMutex{}
)

func SetMonitorStatus(name string, response MonitorResponse) {
	monitorCacheLock.Lock()
	defer monitorCacheLock.Unlock()
	monitorCache[name] = response
}

func GetMonitorStatus(name string) (MonitorResponse, bool) {
	monitorCacheLock.RLock()
	defer monitorCacheLock.RUnlock()
	resp, exists := monitorCache[name]
	return resp, exists
}

func GetAllMonitorStatus() []MonitorResponse {
	monitorCacheLock.RLock()
	defer monitorCacheLock.RUnlock()

	statuses := make([]MonitorResponse, 0, len(monitorCache))
	for _, status := range monitorCache {
		statuses = append(statuses, status)
	}

	return statuses
}
