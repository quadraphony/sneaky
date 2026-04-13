package core_test

import (
	"context"
	"errors"
	"testing"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/core"
	"sneaky-core/internal/runtime"
)

type fakeHandle struct {
	stopErr error
	stopped bool
}

func (h *fakeHandle) Stop(context.Context) error {
	if h.stopErr != nil {
		return h.stopErr
	}
	h.stopped = true
	return nil
}

func (h *fakeHandle) State() runtime.State {
	if h.stopped {
		return runtime.StateStopped
	}
	return runtime.StateRunning
}

type fakeAdapter struct {
	id          string
	validateErr error
	startErr    error
	handle      runtime.Handle
}

func (a *fakeAdapter) Identity() string {
	return a.id
}

func (a *fakeAdapter) Capabilities() adapter.Capabilities {
	return adapter.Capabilities{}
}

func (a *fakeAdapter) ValidateConfig(adapter.StartRequest) error {
	return a.validateErr
}

func (a *fakeAdapter) Start(context.Context, adapter.StartRequest) (runtime.Handle, error) {
	if a.startErr != nil {
		return nil, a.startErr
	}
	if a.handle == nil {
		a.handle = &fakeHandle{}
	}
	return a.handle, nil
}

func TestManagerStartStopLifecycle(t *testing.T) {
	reg := adapter.NewRegistry()
	fake := &fakeAdapter{id: "singbox"}
	if err := reg.Register(fake); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	if err := manager.Start(context.Background(), core.StartRequest{
		AdapterID: "singbox",
		Config: adapter.StartRequest{
			RawConfig: []byte(`{"outbounds":[]}`),
		},
	}); err != nil {
		t.Fatalf("start: %v", err)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateRunning {
		t.Fatalf("expected running state, got %q", snap.State)
	}
	if snap.AdapterID != "singbox" {
		t.Fatalf("expected adapter id singbox, got %q", snap.AdapterID)
	}
	if snap.StartedAt.IsZero() {
		t.Fatal("expected non-zero start time")
	}
	if !snap.Active {
		t.Fatal("expected active snapshot while running")
	}

	statsSnap := manager.Stats()
	if statsSnap.State != runtime.StateRunning {
		t.Fatalf("expected running stats state, got %q", statsSnap.State)
	}
	if statsSnap.SessionsStarted != 1 {
		t.Fatalf("expected 1 started session, got %d", statsSnap.SessionsStarted)
	}
	if len(manager.Logs()) < 2 {
		t.Fatal("expected start logs to be recorded")
	}

	err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"})
	if err == nil {
		t.Fatal("expected second start to fail")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
		t.Fatalf("expected invalid state error, got %v", err)
	}

	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}

	snap = manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state, got %q", snap.State)
	}
	if snap.Active {
		t.Fatal("expected inactive snapshot after stop")
	}

	logs := manager.Logs()
	if len(logs) < 4 {
		t.Fatalf("expected stop logs to be recorded, got %d", len(logs))
	}
	if logs[len(logs)-1].Event != "manager.stop.succeeded" {
		t.Fatalf("expected final log to be stop success, got %q", logs[len(logs)-1].Event)
	}
}

func TestManagerStartValidationFailureKeepsStoppedState(t *testing.T) {
	reg := adapter.NewRegistry()
	if err := reg.Register(&fakeAdapter{
		id:          "singbox",
		validateErr: errors.New("bad config"),
	}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	err := manager.Start(context.Background(), core.StartRequest{
		AdapterID: "singbox",
	})
	if err == nil {
		t.Fatal("expected start validation failure")
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state, got %q", snap.State)
	}
	if snap.LastError == nil {
		t.Fatal("expected last error to be retained")
	}

	statsSnap := manager.Stats()
	if statsSnap.StartFailures != 1 {
		t.Fatalf("expected 1 start failure, got %d", statsSnap.StartFailures)
	}
	logs := manager.Logs()
	if len(logs) == 0 || logs[len(logs)-1].Event != "manager.start.validation_failed" {
		t.Fatalf("expected validation failure log, got %#v", logs)
	}
}

func TestManagerStopFailureKeepsSessionRunning(t *testing.T) {
	reg := adapter.NewRegistry()
	handle := &fakeHandle{stopErr: errors.New("stop failed")}
	if err := reg.Register(&fakeAdapter{
		id:     "singbox",
		handle: handle,
	}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	if err := manager.Start(context.Background(), core.StartRequest{
		AdapterID: "singbox",
	}); err != nil {
		t.Fatalf("start: %v", err)
	}

	err := manager.Stop(context.Background())
	if err == nil {
		t.Fatal("expected stop to fail")
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateRunning {
		t.Fatalf("expected running state after failed stop, got %q", snap.State)
	}
	if snap.LastError == nil {
		t.Fatal("expected stop failure to be stored")
	}

	statsSnap := manager.Stats()
	if statsSnap.StopFailures != 1 {
		t.Fatalf("expected 1 stop failure, got %d", statsSnap.StopFailures)
	}
	logs := manager.Logs()
	if len(logs) == 0 || logs[len(logs)-1].Event != "manager.stop.failed" {
		t.Fatalf("expected stop failure log, got %#v", logs)
	}
}
