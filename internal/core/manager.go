package core

import (
	"context"
	"sync"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/runtime"
)

type StartRequest struct {
	AdapterID string
	Config    adapter.StartRequest
}

// Manager owns the adapter registry and enforces lifecycle rules.
type Manager struct {
	mu       sync.Mutex
	registry *adapter.Registry
	session  *runtime.Session
	state    runtime.State
	lastErr  *Error
}

func NewManager(registry *adapter.Registry) *Manager {
	if registry == nil {
		registry = adapter.NewRegistry()
	}

	return &Manager{
		registry: registry,
		state:    runtime.StateStopped,
	}
}

func (m *Manager) Start(ctx context.Context, req StartRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state != runtime.StateStopped {
		return newError(ErrCodeInvalidState, "core.Manager.Start", "start is only allowed from stopped state", nil)
	}
	if req.AdapterID == "" {
		return newError(ErrCodeInvalidArgument, "core.Manager.Start", "adapter identity is required", nil)
	}

	a, err := m.registry.Resolve(req.AdapterID)
	if err != nil {
		m.lastErr = asCoreError(err, ErrCodeAdapterNotFound, "core.Manager.Start", "failed to resolve adapter")
		return m.lastErr
	}
	if err := a.ValidateConfig(req.Config); err != nil {
		m.lastErr = asCoreError(err, ErrCodeInvalidArgument, "core.Manager.Start", "adapter rejected startup config")
		return m.lastErr
	}

	m.state = runtime.StateStarting
	m.lastErr = nil

	handle, err := a.Start(ctx, req.Config)
	if err != nil {
		m.state = runtime.StateStopped
		m.lastErr = asCoreError(err, ErrCodeStartFailed, "core.Manager.Start", "adapter start failed")
		return m.lastErr
	}

	m.session = &runtime.Session{
		Context: runtime.Context{
			AdapterID: req.AdapterID,
			StartedAt: time.Now().UTC(),
		},
		Handle: handle,
	}
	m.state = runtime.StateRunning
	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state != runtime.StateRunning || m.session == nil || m.session.Handle == nil {
		return newError(ErrCodeInvalidState, "core.Manager.Stop", "stop is only allowed while running", nil)
	}

	m.state = runtime.StateStopping
	if err := m.session.Handle.Stop(ctx); err != nil {
		m.state = runtime.StateRunning
		m.lastErr = asCoreError(err, ErrCodeStopFailed, "core.Manager.Stop", "adapter stop failed")
		return m.lastErr
	}

	m.session = nil
	m.state = runtime.StateStopped
	m.lastErr = nil
	return nil
}

func (m *Manager) Snapshot() Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	snap := Snapshot{
		State:     m.state,
		LastError: m.lastErr,
		Active:    m.state.IsActive(),
	}
	if m.session != nil {
		snap.AdapterID = m.session.Context.AdapterID
		snap.StartedAt = m.session.Context.StartedAt
	}
	return snap
}

func (m *Manager) Registry() *adapter.Registry {
	return m.registry
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
