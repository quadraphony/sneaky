package stats

import (
	"testing"
	"time"

	"sneaky-core/internal/runtime"
)

func TestTrackerSnapshotsLifecycleCounters(t *testing.T) {
	tracker := NewTracker()

	tracker.RecordStarting("singbox")
	startedAt := time.Now().UTC().Add(-2 * time.Second)
	tracker.RecordRunning("singbox", startedAt)

	snap := tracker.Snapshot()
	if snap.State != runtime.StateRunning {
		t.Fatalf("expected running state, got %q", snap.State)
	}
	if snap.AdapterID != "singbox" {
		t.Fatalf("expected adapter singbox, got %q", snap.AdapterID)
	}
	if snap.SessionsStarted != 1 {
		t.Fatalf("expected 1 session started, got %d", snap.SessionsStarted)
	}
	if snap.Uptime <= 0 {
		t.Fatalf("expected positive uptime, got %s", snap.Uptime)
	}

	tracker.RecordStopping()
	tracker.RecordStopFailure()
	snap = tracker.Snapshot()
	if snap.StopFailures != 1 {
		t.Fatalf("expected 1 stop failure, got %d", snap.StopFailures)
	}

	tracker.RecordStopped()
	tracker.RecordStartFailure()
	snap = tracker.Snapshot()
	if snap.StartFailures != 1 {
		t.Fatalf("expected 1 start failure, got %d", snap.StartFailures)
	}
	if snap.State != runtime.StateStopped {
		t.Fatalf("expected stopped state, got %q", snap.State)
	}
}
