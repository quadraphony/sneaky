package sneaky

import (
	"context"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/core"
	"sneaky-core/internal/runtime"
)

// Manager is the stable public wrapper around the internal core manager.
type Manager struct {
	core *core.Manager
}

type State string

const (
	StateStopped  State = State(runtime.StateStopped)
	StateStarting State = State(runtime.StateStarting)
	StateRunning  State = State(runtime.StateRunning)
	StateStopping State = State(runtime.StateStopping)
)

type Snapshot struct {
	State     State
	AdapterID string
	StartedAt time.Time
	Active    bool
	LastError error
}

type StartRequest struct {
	AdapterID  string
	ConfigPath string
	RawConfig  []byte
}

func New() *Manager {
	return &Manager{
		core: core.NewManager(adapter.NewRegistry()),
	}
}

func (m *Manager) Start(ctx context.Context, req StartRequest) error {
	return m.core.Start(ctx, core.StartRequest{
		AdapterID: req.AdapterID,
		Config: adapter.StartRequest{
			ConfigPath: req.ConfigPath,
			RawConfig:  req.RawConfig,
		},
	})
}

func (m *Manager) Stop(ctx context.Context) error {
	return m.core.Stop(ctx)
}

func (m *Manager) Snapshot() Snapshot {
	snap := m.core.Snapshot()
	return Snapshot{
		State:     State(snap.State),
		AdapterID: snap.AdapterID,
		StartedAt: snap.StartedAt,
		Active:    snap.Active,
		LastError: snap.LastError,
	}
}
