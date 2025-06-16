package register

import (
	"sync"
	"testing"
	"time"
	"wavezync/pulse-bridge/internal/config"
)

func mockMonitor(name, interval, timeout, mtype string) config.Monitor {
	return config.Monitor{
		Name:     name,
		Type:     mtype,
		Interval: interval,
		Timeout:  timeout,
	}
}

func TestRunMonitorWorker_TimerBehavior(t *testing.T) {
	interval := "100ms"
	timeout := "200ms"
	m := mockMonitor("test-monitor", interval, timeout, "http")

	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	calls := 0

	mockTimer := func(mntr *config.Monitor) {
		calls++
		if calls == 2 {
			elapsed := time.Since(start)
			if elapsed < 90*time.Millisecond || elapsed > 200*time.Millisecond {
				t.Errorf("Second call timing off: got %v, want between 90ms and 200ms", elapsed)
			}
			wg.Done()
		}
	}

	go runMonitorWorker(&m, mockTimer)
	wg.Wait()
}
