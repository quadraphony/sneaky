package sneaky

import (
	"context"
	"testing"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/runtime"
)

type mockHandle struct {
	stopped bool
}

func (h *mockHandle) Stop(ctx context.Context) error {
	h.stopped = true
	return nil
}

func (h *mockHandle) State() runtime.State {
	if h.stopped {
		return runtime.StateStopped
	}
	return runtime.StateRunning
}

type mockAdapter struct{}

func (a *mockAdapter) Identity() string {
	return "mock"
}

func (a *mockAdapter) Capabilities() adapter.Capabilities {
	return adapter.Capabilities{}
}

func (a *mockAdapter) ValidateConfig(req adapter.StartRequest) error {
	if req.ConfigPath == "invalid.json" {
		return context.DeadlineExceeded // just some error
	}
	return nil
}

func (a *mockAdapter) Start(ctx context.Context, req adapter.StartRequest) (runtime.Handle, error) {
	if req.ConfigPath == "fail_start.json" {
		return nil, context.DeadlineExceeded
	}
	return &mockHandle{}, nil
}

func TestManagerLifecycleWithMock(t *testing.T) {
	registry := adapter.NewRegistry()
	registry.MustRegister(&mockAdapter{})
	manager := NewWithRegistry(registry)

	// Test Start
	err := manager.Start(context.Background(), StartRequest{
		AdapterID:  "mock",
		ConfigPath: "fake.json",
	})
	if err != nil {
		t.Fatalf("failed to start mock adapter: %v", err)
	}

	snap := manager.Snapshot()
	if snap.State != StateRunning {
		t.Fatalf("expected running state, got %v", snap.State)
	}
	if snap.AdapterID != "mock" {
		t.Fatalf("expected adapter mock, got %v", snap.AdapterID)
	}

	// Test Stats
	stats := manager.Stats()
	if stats.SessionsStarted != 1 {
		t.Fatalf("expected 1 session started, got %d", stats.SessionsStarted)
	}
	if stats.Uptime <= 0 {
		t.Fatalf("expected positive uptime, got %v", stats.Uptime)
	}

	// Test Stop
	stopCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = manager.Stop(stopCtx)
	if err != nil {
		t.Fatalf("failed to stop mock adapter: %v", err)
	}

	snap = manager.Snapshot()
	if snap.State != StateStopped {
		t.Fatalf("expected stopped state, got %v", snap.State)
	}
}

func TestManagerStartValidationFailure(t *testing.T) {
	registry := adapter.NewRegistry()
	registry.MustRegister(&mockAdapter{})
	manager := NewWithRegistry(registry)

	err := manager.Start(context.Background(), StartRequest{
		AdapterID:  "mock",
		ConfigPath: "invalid.json",
	})
	if err == nil {
		t.Fatal("expected start to fail due to validation")
	}

	snap := manager.Snapshot()
	if snap.State != StateStopped {
		t.Fatalf("expected stopped state, got %v", snap.State)
	}
	if snap.LastError == nil {
		t.Fatal("expected last error to be set")
	}
}

func TestManagerStartFailure(t *testing.T) {
	registry := adapter.NewRegistry()
	registry.MustRegister(&mockAdapter{})
	manager := NewWithRegistry(registry)

	err := manager.Start(context.Background(), StartRequest{
		AdapterID:  "mock",
		ConfigPath: "fail_start.json",
	})
	if err == nil {
		t.Fatal("expected start to fail")
	}

	stats := manager.Stats()
	if stats.StartFailures != 1 {
		t.Fatalf("expected 1 start failure, got %d", stats.StartFailures)
	}
}
