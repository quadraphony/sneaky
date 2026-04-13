package singbox

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"sneaky-core/internal/adapter"
)

const validConfig = `{
  "log": {"level": "warn"},
  "outbounds": [
    {"type": "direct", "tag": "direct"}
  ]
}`

func requireSingboxBinary(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("sing-box"); err != nil {
		t.Skip("sing-box binary not available")
	}
}

func TestValidateConfigWithInlineConfig(t *testing.T) {
	requireSingboxBinary(t)

	a := New("")
	err := a.ValidateConfig(adapter.StartRequest{
		RawConfig: []byte(validConfig),
	})
	if err != nil {
		t.Fatalf("validate inline config: %v", err)
	}
}

func TestValidateConfigWithPath(t *testing.T) {
	requireSingboxBinary(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(validConfig), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	a := New("")
	err := a.ValidateConfig(adapter.StartRequest{
		ConfigPath: path,
	})
	if err != nil {
		t.Fatalf("validate file config: %v", err)
	}
}

func TestStartAndStop(t *testing.T) {
	requireSingboxBinary(t)

	a := New("")
	handle, err := a.Start(context.Background(), adapter.StartRequest{
		RawConfig: []byte(validConfig),
	})
	if err != nil {
		t.Fatalf("start adapter: %v", err)
	}
	if handle.State().String() != "running" {
		t.Fatalf("expected running state, got %q", handle.State())
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := handle.Stop(stopCtx); err != nil {
		t.Fatalf("stop adapter: %v", err)
	}
	if handle.State().String() != "stopped" {
		t.Fatalf("expected stopped state, got %q", handle.State())
	}
}
