package mock

import "testing"

func TestMockStartStop(t *testing.T) {
	m := New()
	if err := m.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	st := m.Status()
	if running, ok := st["running"]; !ok || running != true {
		t.Fatalf("expected running=true, got %v", st)
	}
	if err := m.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}
