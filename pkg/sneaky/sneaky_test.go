package sneaky

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

const singboxConfig = `{
  "log": {"level": "warn"},
  "outbounds": [
    {"type": "direct", "tag": "direct"}
  ]
}`

func TestManagerRepeatedStartStop(t *testing.T) {
	if _, err := exec.LookPath("sing-box"); err != nil {
		t.Skip("sing-box binary not available")
	}

	manager := New()
	for i := 0; i < 2; i++ {
		if err := manager.Start(context.Background(), StartRequest{
			AdapterID: "singbox",
			RawConfig: []byte(singboxConfig),
		}); err != nil {
			t.Fatalf("start iteration %d: %v", i, err)
		}

		stats := manager.Stats()
		if stats.State != StateRunning {
			t.Fatalf("expected running stats state, got %q", stats.State)
		}

		stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		if err := manager.Stop(stopCtx); err != nil {
			cancel()
			t.Fatalf("stop iteration %d: %v", i, err)
		}
		cancel()
	}

	stats := manager.Stats()
	if stats.SessionsStarted != 2 {
		t.Fatalf("expected 2 sessions started, got %d", stats.SessionsStarted)
	}
	if len(manager.Logs()) == 0 {
		t.Fatal("expected logs to be recorded")
	}
}
