package core_test

import (
	"context"
	"errors"
	"testing"
	"time"

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
	startGate   <-chan struct{}
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
	if a.startGate != nil {
		<-a.startGate
	}
	if a.handle == nil || a.handle.State() == runtime.StateStopped {
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

func TestManagerCanRestartAfterCleanStop(t *testing.T) {
	reg := adapter.NewRegistry()
	if err := reg.Register(&fakeAdapter{id: "singbox"}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	for i := 0; i < 2; i++ {
		if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err != nil {
			t.Fatalf("start iteration %d: %v", i, err)
		}
		if err := manager.Stop(context.Background()); err != nil {
			t.Fatalf("stop iteration %d: %v", i, err)
		}
	}

	statsSnap := manager.Stats()
	if statsSnap.SessionsStarted != 2 {
		t.Fatalf("expected 2 sessions started, got %d", statsSnap.SessionsStarted)
	}
}

func TestManagerReconcilesUnexpectedRuntimeExit(t *testing.T) {
	reg := adapter.NewRegistry()
	handle := &fakeHandle{}
	if err := reg.Register(&fakeAdapter{
		id:     "singbox",
		handle: handle,
	}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err != nil {
		t.Fatalf("start: %v", err)
	}

	handle.stopped = true

	snap := manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state after runtime exit, got %q", snap.State)
	}
	if snap.LastError == nil || snap.LastError.Code != core.ErrCodeRuntimeExited {
		t.Fatalf("expected runtime exited error, got %#v", snap.LastError)
	}
	if snap.Active {
		t.Fatal("expected inactive snapshot after runtime exit")
	}

	statsSnap := manager.Stats()
	if statsSnap.State != runtime.StateStopped {
		t.Fatalf("expected stopped stats state, got %q", statsSnap.State)
	}

	if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err != nil {
		t.Fatalf("restart after runtime exit: %v", err)
	}
}

func TestManagerRejectsConcurrentStartWhileStarting(t *testing.T) {
	reg := adapter.NewRegistry()
	startGate := make(chan struct{})
	if err := reg.Register(&fakeAdapter{
		id:        "singbox",
		startGate: startGate,
	}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)

	startResult := make(chan error, 1)
	go func() {
		startResult <- manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snap := manager.Snapshot()
		if snap.State == runtime.StateStarting {
			if snap.AdapterID != "singbox" {
				t.Fatalf("expected pending adapter id singbox, got %q", snap.AdapterID)
			}
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStarting {
		t.Fatalf("expected starting state, got %q", snap.State)
	}

	secondStart := make(chan error, 1)
	go func() {
		secondStart <- manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"})
	}()

	select {
	case err := <-secondStart:
		var coreErr *core.Error
		if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
			t.Fatalf("expected invalid state error for concurrent start, got %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("concurrent start blocked instead of failing immediately")
	}

	close(startGate)

	if err := <-startResult; err != nil {
		t.Fatalf("first start failed: %v", err)
	}
	if manager.Snapshot().State != runtime.StateRunning {
		t.Fatalf("expected running state after first start completes, got %q", manager.Snapshot().State)
	}
}

func TestManagerStopWhileStartingFailsWithoutMutation(t *testing.T) {
	reg := adapter.NewRegistry()
	startGate := make(chan struct{})
	if err := reg.Register(&fakeAdapter{
		id:        "singbox",
		startGate: startGate,
	}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)

	startResult := make(chan error, 1)
	go func() {
		startResult <- manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snap := manager.Snapshot()
		if snap.State == runtime.StateStarting {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStarting {
		t.Fatalf("expected starting state, got %q", snap.State)
	}

	err := manager.Stop(context.Background())
	if err == nil {
		t.Fatal("expected stop while starting to fail")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
		t.Fatalf("expected invalid state error, got %v", err)
	}

	snap = manager.Snapshot()
	if snap.State != runtime.StateStarting {
		t.Fatalf("expected manager to remain in starting state, got %q", snap.State)
	}
	if snap.AdapterID != "singbox" {
		t.Fatalf("expected pending adapter id singbox, got %q", snap.AdapterID)
	}
	if snap.LastError != nil {
		t.Fatalf("expected no hidden last error mutation, got %#v", snap.LastError)
	}

	close(startGate)
	if err := <-startResult; err != nil {
		t.Fatalf("first start failed: %v", err)
	}
}

func TestManagerStopWhileIdleDoesNotMutateState(t *testing.T) {
	manager := core.NewManager(adapter.NewRegistry())

	err := manager.Stop(context.Background())
	if err == nil {
		t.Fatal("expected stop while idle to fail")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
		t.Fatalf("expected invalid state error, got %v", err)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state, got %q", snap.State)
	}
	if snap.LastError != nil {
		t.Fatalf("expected no last error mutation, got %#v", snap.LastError)
	}
}

func TestManagerStopAfterRuntimeExitPreservesExitReason(t *testing.T) {
	reg := adapter.NewRegistry()
	handle := &fakeHandle{}
	if err := reg.Register(&fakeAdapter{
		id:     "singbox",
		handle: handle,
	}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err != nil {
		t.Fatalf("start: %v", err)
	}

	handle.stopped = true

	err := manager.Stop(context.Background())
	if err == nil {
		t.Fatal("expected stop after runtime exit to fail")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
		t.Fatalf("expected invalid state error, got %v", err)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state, got %q", snap.State)
	}
	if snap.LastError == nil || snap.LastError.Code != core.ErrCodeRuntimeExited {
		t.Fatalf("expected preserved runtime exited error, got %#v", snap.LastError)
	}
}

func TestManagerStopTwiceReturnsPredictableError(t *testing.T) {
	reg := adapter.NewRegistry()
	if err := reg.Register(&fakeAdapter{id: "singbox"}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("first stop: %v", err)
	}

	err := manager.Stop(context.Background())
	if err == nil {
		t.Fatal("expected second stop to fail")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
		t.Fatalf("expected invalid state error, got %v", err)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state after second stop, got %q", snap.State)
	}
	if snap.LastError != nil {
		t.Fatalf("expected no hidden last error mutation after second stop, got %#v", snap.LastError)
	}
}

func TestManagerStartFailureClearsPendingStateAndAllowsRetry(t *testing.T) {
	reg := adapter.NewRegistry()
	fake := &fakeAdapter{
		id:       "singbox",
		startErr: errors.New("spawn failed"),
	}
	if err := reg.Register(fake); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)

	err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"})
	if err == nil {
		t.Fatal("expected start failure")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeStartFailed {
		t.Fatalf("expected start_failed error, got %v", err)
	}

	snap := manager.Snapshot()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state after failed start, got %q", snap.State)
	}
	if snap.AdapterID != "" {
		t.Fatalf("expected no active adapter after failed start, got %q", snap.AdapterID)
	}
	if snap.LastError == nil || snap.LastError.Code != core.ErrCodeStartFailed {
		t.Fatalf("expected retained start failure, got %#v", snap.LastError)
	}

	statsSnap := manager.Stats()
	if statsSnap.State != runtime.StateStopped {
		t.Fatalf("expected stopped stats state after failed start, got %q", statsSnap.State)
	}
	if statsSnap.StartFailures != 1 {
		t.Fatalf("expected 1 start failure, got %d", statsSnap.StartFailures)
	}
	if statsSnap.AdapterID != "" {
		t.Fatalf("expected cleared adapter id after failed start, got %q", statsSnap.AdapterID)
	}

	fake.startErr = nil
	if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err != nil {
		t.Fatalf("retry start after failure: %v", err)
	}
	if manager.Snapshot().State != runtime.StateRunning {
		t.Fatalf("expected running state after retry, got %q", manager.Snapshot().State)
	}
}

func TestManagerStopAfterFailedStartPreservesFailureReason(t *testing.T) {
	reg := adapter.NewRegistry()
	fake := &fakeAdapter{
		id:       "singbox",
		startErr: errors.New("spawn failed"),
	}
	if err := reg.Register(fake); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	manager := core.NewManager(reg)
	if err := manager.Start(context.Background(), core.StartRequest{AdapterID: "singbox"}); err == nil {
		t.Fatal("expected start failure")
	}

	err := manager.Stop(context.Background())
	if err == nil {
		t.Fatal("expected stop after failed start to fail")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeInvalidState {
		t.Fatalf("expected invalid state error, got %v", err)
	}

	snap := manager.Snapshot()
	if snap.LastError == nil || snap.LastError.Code != core.ErrCodeStartFailed {
		t.Fatalf("expected preserved start failure, got %#v", snap.LastError)
	}
}

func TestManagerResolveFailureResetsStartingStats(t *testing.T) {
	manager := core.NewManager(adapter.NewRegistry())

	err := manager.Start(context.Background(), core.StartRequest{AdapterID: "missing"})
	if err == nil {
		t.Fatal("expected resolve failure")
	}
	var coreErr *core.Error
	if !errors.As(err, &coreErr) || coreErr.Code != core.ErrCodeAdapterNotFound {
		t.Fatalf("expected adapter_not_found error, got %v", err)
	}

	snap := manager.Stats()
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped stats state after resolve failure, got %q", snap.State)
	}
	if snap.StartFailures != 1 {
		t.Fatalf("expected 1 start failure after resolve failure, got %d", snap.StartFailures)
	}
	if snap.AdapterID != "" {
		t.Fatalf("expected cleared adapter id after resolve failure, got %q", snap.AdapterID)
	}
}
