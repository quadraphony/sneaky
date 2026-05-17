package sneaky

import (
	"context"
	"os"
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
			AdapterID: AdapterIDSingbox,
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

	for _, entry := range manager.Logs() {
		if entry.Message == "" {
			t.Error("expected non-empty log message")
		}
	}
}

func TestInspectConfigPath(t *testing.T) {
	if _, err := exec.LookPath("sing-box"); err != nil {
		t.Skip("sing-box binary not available")
	}

	tempDir := t.TempDir()
	configPath := tempDir + "/config.json"
	err := os.WriteFile(configPath, []byte(singboxConfig), 0644)
	if err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	metadata, err := InspectConfigPath(configPath)
	if err != nil {
		t.Fatalf("InspectConfigPath failed: %v", err)
	}

	if metadata.AdapterID != AdapterIDSingbox {
		t.Errorf("expected adapter %v, got %v", AdapterIDSingbox, metadata.AdapterID)
	}
}

func TestPublicStateHelpers(t *testing.T) {
	if !StateStopped.CanStart() {
		t.Fatal("expected stopped state to allow start")
	}
	if StateRunning.CanStart() {
		t.Fatal("expected running state to reject start")
	}
	if !StateRunning.CanStop() {
		t.Fatal("expected running state to allow stop")
	}
	if StateStarting.CanStop() {
		t.Fatal("expected starting state to reject stop")
	}
	if !StateStarting.IsActive() {
		t.Fatal("expected starting state to be active")
	}
	if StateStopped.IsActive() {
		t.Fatal("expected stopped state to be inactive")
	}
	if StateStopping.String() != "stopping" {
		t.Fatalf("unexpected state string: %q", StateStopping.String())
	}
}
