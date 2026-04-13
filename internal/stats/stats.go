package stats

import (
	"sync"
	"time"

	"sneaky-core/internal/runtime"
)

// Tracker records manager-level runtime counters and timing.
type Tracker struct {
	mu               sync.Mutex
	state            runtime.State
	adapterID        string
	startedAt        time.Time
	lastTransitionAt time.Time
	sessionsStarted  uint64
	startFailures    uint64
	stopFailures     uint64
}

func NewTracker() *Tracker {
	now := time.Now().UTC()
	return &Tracker{
		state:            runtime.StateStopped,
		lastTransitionAt: now,
	}
}

func (t *Tracker) RecordStarting(adapterID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state = runtime.StateStarting
	t.adapterID = adapterID
	t.lastTransitionAt = time.Now().UTC()
}

func (t *Tracker) RecordRunning(adapterID string, startedAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state = runtime.StateRunning
	t.adapterID = adapterID
	t.startedAt = startedAt.UTC()
	t.lastTransitionAt = t.startedAt
	t.sessionsStarted++
}

func (t *Tracker) RecordStartFailure() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state = runtime.StateStopped
	t.adapterID = ""
	t.startedAt = time.Time{}
	t.lastTransitionAt = time.Now().UTC()
	t.startFailures++
}

func (t *Tracker) RecordStopping() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state = runtime.StateStopping
	t.lastTransitionAt = time.Now().UTC()
}

func (t *Tracker) RecordStopped() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state = runtime.StateStopped
	t.adapterID = ""
	t.startedAt = time.Time{}
	t.lastTransitionAt = time.Now().UTC()
}

func (t *Tracker) RecordStopFailure() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state = runtime.StateRunning
	t.lastTransitionAt = time.Now().UTC()
	t.stopFailures++
}

func (t *Tracker) Snapshot() Snapshot {
	t.mu.Lock()
	defer t.mu.Unlock()

	snapshot := Snapshot{
		State:            t.state,
		AdapterID:        t.adapterID,
		StartedAt:        t.startedAt,
		LastTransitionAt: t.lastTransitionAt,
		SessionsStarted:  t.sessionsStarted,
		StartFailures:    t.startFailures,
		StopFailures:     t.stopFailures,
	}
	if t.state == runtime.StateRunning && !t.startedAt.IsZero() {
		snapshot.Uptime = time.Since(t.startedAt)
	}
	return snapshot
}
