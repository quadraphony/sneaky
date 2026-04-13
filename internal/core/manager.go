package core

import (
	"context"
	"sync"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/logx"
	"sneaky-core/internal/runtime"
	"sneaky-core/internal/stats"
)

type StartRequest struct {
	AdapterID string
	Config    adapter.StartRequest
}

// Manager owns the adapter registry and enforces lifecycle rules.
type Manager struct {
	mu               sync.Mutex
	registry         *adapter.Registry
	logger           *logx.Logger
	stats            *stats.Tracker
	session          *runtime.Session
	state            runtime.State
	lastErr          *Error
	pendingAdapterID string
}

func NewManager(registry *adapter.Registry) *Manager {
	if registry == nil {
		registry = adapter.NewRegistry()
	}

	return &Manager{
		registry: registry,
		logger:   logx.New(256),
		stats:    stats.NewTracker(),
		state:    runtime.StateStopped,
	}
}

func (m *Manager) Start(ctx context.Context, req StartRequest) error {
	m.mu.Lock()

	m.reconcileSessionLocked()

	if !m.state.CanStart() {
		err := newError(ErrCodeInvalidState, "core.Manager.Start", "start is only allowed from stopped state", nil)
		m.lastErr = err
		m.logError("manager.start.rejected", err, map[string]string{"state": m.state.String()})
		m.mu.Unlock()
		return err
	}
	if req.AdapterID == "" {
		err := newError(ErrCodeInvalidArgument, "core.Manager.Start", "adapter identity is required", nil)
		m.lastErr = err
		m.logError("manager.start.rejected", err, nil)
		m.mu.Unlock()
		return err
	}

	m.logger.Info("manager.start.requested", "manager start requested", map[string]string{"adapter_id": req.AdapterID})
	m.state = runtime.StateStarting
	m.pendingAdapterID = req.AdapterID
	m.stats.RecordStarting(req.AdapterID)
	m.lastErr = nil
	m.mu.Unlock()

	a, err := m.registry.Resolve(req.AdapterID)
	if err != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.clearPendingStartLocked()
		m.stats.RecordStartFailure()
		m.lastErr = asCoreError(err, ErrCodeAdapterNotFound, "core.Manager.Start", "failed to resolve adapter")
		m.logError("manager.start.resolve_failed", m.lastErr, map[string]string{"adapter_id": req.AdapterID})
		return m.lastErr
	}
	if err := a.ValidateConfig(req.Config); err != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.clearPendingStartLocked()
		m.lastErr = asCoreError(err, ErrCodeInvalidArgument, "core.Manager.Start", "adapter rejected startup config")
		m.stats.RecordStartFailure()
		m.logError("manager.start.validation_failed", m.lastErr, map[string]string{"adapter_id": req.AdapterID})
		return m.lastErr
	}

	handle, err := a.Start(ctx, req.Config)
	if err != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.clearPendingStartLocked()
		m.stats.RecordStartFailure()
		m.lastErr = asCoreError(err, ErrCodeStartFailed, "core.Manager.Start", "adapter start failed")
		m.logError("manager.start.failed", m.lastErr, map[string]string{"adapter_id": req.AdapterID})
		return m.lastErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	startedAt := time.Now().UTC()
	m.session = &runtime.Session{
		Context: runtime.Context{
			AdapterID: req.AdapterID,
			StartedAt: startedAt,
		},
		Handle: handle,
	}
	m.pendingAdapterID = ""
	m.state = runtime.StateRunning
	m.stats.RecordRunning(req.AdapterID, startedAt)
	m.logger.Info("manager.start.succeeded", "manager start succeeded", map[string]string{"adapter_id": req.AdapterID})
	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reconcileSessionLocked()

	if !m.state.CanStop() || m.session == nil || m.session.Handle == nil {
		err := newError(ErrCodeInvalidState, "core.Manager.Stop", "stop is only allowed while running", nil)
		m.logError("manager.stop.rejected", err, map[string]string{"state": m.state.String()})
		return err
	}

	adapterID := m.session.Context.AdapterID
	m.logger.Info("manager.stop.requested", "manager stop requested", map[string]string{"adapter_id": adapterID})
	m.state = runtime.StateStopping
	m.stats.RecordStopping()
	if err := m.session.Handle.Stop(ctx); err != nil {
		m.state = runtime.StateRunning
		m.stats.RecordStopFailure()
		m.lastErr = asCoreError(err, ErrCodeStopFailed, "core.Manager.Stop", "adapter stop failed")
		m.logError("manager.stop.failed", m.lastErr, map[string]string{"adapter_id": adapterID})
		return m.lastErr
	}

	m.session = nil
	m.state = runtime.StateStopped
	m.stats.RecordStopped()
	m.lastErr = nil
	m.logger.Info("manager.stop.succeeded", "manager stop succeeded", map[string]string{"adapter_id": adapterID})
	return nil
}

func (m *Manager) Snapshot() Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reconcileSessionLocked()

	snap := Snapshot{
		State:     m.state,
		LastError: m.lastErr,
		Active:    m.state.IsActive(),
	}
	if m.session != nil {
		snap.AdapterID = m.session.Context.AdapterID
		snap.StartedAt = m.session.Context.StartedAt
	} else if m.pendingAdapterID != "" {
		snap.AdapterID = m.pendingAdapterID
	}
	return snap
}

func (m *Manager) Registry() *adapter.Registry {
	return m.registry
}

func (m *Manager) Logs() []logx.Entry {
	return m.logger.Entries()
}

func (m *Manager) Stats() stats.Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reconcileSessionLocked()

	snapshot := m.stats.Snapshot()
	snapshot.State = m.state
	if m.session != nil {
		snapshot.AdapterID = m.session.Context.AdapterID
		snapshot.StartedAt = m.session.Context.StartedAt
	} else if m.pendingAdapterID != "" {
		snapshot.AdapterID = m.pendingAdapterID
	}
	if snapshot.State == runtime.StateRunning && !snapshot.StartedAt.IsZero() {
		snapshot.Uptime = time.Since(snapshot.StartedAt)
	}
	return snapshot
}

func (m *Manager) logError(event string, err error, fields map[string]string) {
	if err == nil {
		return
	}

	merged := make(map[string]string, len(fields)+1)
	for key, value := range fields {
		merged[key] = value
	}
	merged["error"] = err.Error()
	m.logger.Error(event, "manager operation failed", merged)
}

func (m *Manager) reconcileSessionLocked() {
	if m.session == nil || m.session.Handle == nil {
		return
	}
	if m.state != runtime.StateRunning && m.state != runtime.StateStopping {
		return
	}
	if m.session.Handle.State() != runtime.StateStopped {
		return
	}

	adapterID := m.session.Context.AdapterID
	m.session = nil
	m.pendingAdapterID = ""
	m.state = runtime.StateStopped
	if m.lastErr == nil {
		m.lastErr = newError(ErrCodeRuntimeExited, "core.Manager", "adapter runtime exited unexpectedly", nil)
	}
	m.stats.RecordStopped()
	m.logError("manager.runtime.exited", m.lastErr, map[string]string{"adapter_id": adapterID})
}

func (m *Manager) clearPendingStartLocked() {
	m.state = runtime.StateStopped
	m.session = nil
	m.pendingAdapterID = ""
}

func asCoreError(err error, fallbackCode ErrorCode, op, message string) *Error {
	if err == nil {
		return nil
	}

	if coreErr, ok := err.(*Error); ok {
		return &Error{
			Code:    coreErr.Code,
			Op:      op,
			Message: message,
			Err:     coreErr,
		}
	}

	return newError(fallbackCode, op, message, err)
}
