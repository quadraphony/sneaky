package mock

import (
	"sync"
	"time"
	"errors"
)

// MockAdapter is a minimal adapter implementation for tests and CI.
// It follows a simple Start/Stop/Status contract.
type MockAdapter struct {
	mu sync.Mutex
	running bool
	startedAt time.Time
}

func New() *MockAdapter {
	return &MockAdapter{}
}

func (m *MockAdapter) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return errors.New("already running")
	}
	m.running = true
	m.startedAt = time.Now()
	return nil
}

func (m *MockAdapter) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return errors.New("not running")
	}
	m.running = false
	return nil
}

func (m *MockAdapter) Status() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return map[string]interface{}{"running": false}
	}
	return map[string]interface{}{
		"running": true,
		"started_at": m.startedAt.Format(time.RFC3339),
	}
}
