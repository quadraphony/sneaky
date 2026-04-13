package sneaky

import (
	"context"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/adapters/singbox"
	sshadapter "sneaky-core/internal/adapters/ssh"
	"sneaky-core/internal/core"
	"sneaky-core/internal/logx"
	"sneaky-core/internal/runtime"
	"sneaky-core/internal/stats"
)

// Manager is the stable public wrapper around the internal core manager.
type Manager struct {
	core *core.Manager
}

type AdapterID string

const (
	AdapterIDSingbox AdapterID = "singbox"
	AdapterIDSSH     AdapterID = "ssh"
)

type State string

const (
	StateStopped  State = State(runtime.StateStopped)
	StateStarting State = State(runtime.StateStarting)
	StateRunning  State = State(runtime.StateRunning)
	StateStopping State = State(runtime.StateStopping)
)

type Snapshot struct {
	State     State
	AdapterID AdapterID
	StartedAt time.Time
	Active    bool
	LastError error
}

type LogEntry struct {
	Time    time.Time
	Level   string
	Event   string
	Message string
	Fields  map[string]string
}

type StatsSnapshot struct {
	State            State
	AdapterID        AdapterID
	StartedAt        time.Time
	LastTransitionAt time.Time
	Uptime           time.Duration
	SessionsStarted  uint64
	StartFailures    uint64
	StopFailures     uint64
}

type StartRequest struct {
	AdapterID  AdapterID
	ConfigPath string
	RawConfig  []byte
}

func New() *Manager {
	registry := adapter.NewRegistry()
	registry.MustRegister(singbox.New(""))
	registry.MustRegister(sshadapter.New(""))

	return &Manager{
		core: core.NewManager(registry),
	}
}

func (s State) String() string {
	return string(s)
}

func (s State) IsActive() bool {
	return runtime.State(s).IsActive()
}

func (s State) CanStart() bool {
	return runtime.State(s).CanStart()
}

func (s State) CanStop() bool {
	return runtime.State(s).CanStop()
}

func (m *Manager) Start(ctx context.Context, req StartRequest) error {
	return m.core.Start(ctx, core.StartRequest{
		AdapterID: string(req.AdapterID),
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
		AdapterID: AdapterID(snap.AdapterID),
		StartedAt: snap.StartedAt,
		Active:    snap.Active,
		LastError: snap.LastError,
	}
}

func (m *Manager) Logs() []LogEntry {
	entries := m.core.Logs()
	out := make([]LogEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, fromLogEntry(entry))
	}
	return out
}

func (m *Manager) Stats() StatsSnapshot {
	return fromStatsSnapshot(m.core.Stats())
}

func fromLogEntry(entry logx.Entry) LogEntry {
	fields := make(map[string]string, len(entry.Fields))
	for key, value := range entry.Fields {
		fields[key] = value
	}

	return LogEntry{
		Time:    entry.Time,
		Level:   string(entry.Level),
		Event:   entry.Event,
		Message: entry.Message,
		Fields:  fields,
	}
}

func fromStatsSnapshot(snapshot stats.Snapshot) StatsSnapshot {
	return StatsSnapshot{
		State:            State(snapshot.State),
		AdapterID:        AdapterID(snapshot.AdapterID),
		StartedAt:        snapshot.StartedAt,
		LastTransitionAt: snapshot.LastTransitionAt,
		Uptime:           snapshot.Uptime,
		SessionsStarted:  snapshot.SessionsStarted,
		StartFailures:    snapshot.StartFailures,
		StopFailures:     snapshot.StopFailures,
	}
}
