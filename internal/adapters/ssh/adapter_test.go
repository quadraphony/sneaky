package ssh

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"sneaky-core/internal/adapter"
)

const validSSHConfig = `{
  "ssh_tunnel": {
    "host": "example.com",
    "user": "demo",
    "local_socks_port": 1080
  }
}`

func TestValidateConfig(t *testing.T) {
	binary := writeFakeSSHBinary(t)
	a := New(binary)

	if err := a.ValidateConfig(adapter.StartRequest{RawConfig: []byte(validSSHConfig)}); err != nil {
		t.Fatalf("validate config: %v", err)
	}
}

func TestStartAndStop(t *testing.T) {
	binary := writeFakeSSHBinary(t)
	a := New(binary)

	handle, err := a.Start(context.Background(), adapter.StartRequest{RawConfig: []byte(validSSHConfig)})
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
}

func writeFakeSSHBinary(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "ssh")
	script := `#!/usr/bin/env bash
set -euo pipefail
if [[ "${1:-}" == "-G" ]]; then
  exit 0
fi
trap 'exit 0' TERM INT
while true; do
  sleep 1
done
`
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake ssh binary: %v", err)
	}
	return path
}

func TestLoadTunnelConfigDefaults(t *testing.T) {
	cfg, err := loadTunnelConfig(adapter.StartRequest{RawConfig: []byte(validSSHConfig)})
	if err != nil {
		t.Fatalf("load tunnel config: %v", err)
	}
	if cfg.Port != 22 {
		t.Fatalf("expected default port 22, got %d", cfg.Port)
	}
	if !strings.EqualFold(cfg.StrictHostKeyChecking, "accept-new") {
		t.Fatalf("expected default strict host key checking, got %q", cfg.StrictHostKeyChecking)
	}
	if cfg.ConnectTimeoutSeconds != 10 {
		t.Fatalf("expected default connect timeout 10, got %d", cfg.ConnectTimeoutSeconds)
	}
	if cfg.ServerAliveInterval != 30 {
		t.Fatalf("expected default server alive interval 30, got %d", cfg.ServerAliveInterval)
	}
	if cfg.ServerAliveCountMax != 3 {
		t.Fatalf("expected default server alive count max 3, got %d", cfg.ServerAliveCountMax)
	}
}

func TestLoadTunnelConfigResolvesKnownHostsFile(t *testing.T) {
	dir := t.TempDir()
	knownHostsPath := filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(knownHostsPath, []byte("example"), 0o600); err != nil {
		t.Fatalf("write known hosts: %v", err)
	}

	cfg, err := loadTunnelConfig(adapter.StartRequest{RawConfig: []byte(`{
	  "ssh_tunnel": {
	    "host": "example.com",
	    "user": "demo",
	    "local_socks_port": 1080,
	    "known_hosts_file": "` + knownHostsPath + `"
	  }
	}`)})
	if err != nil {
		t.Fatalf("load tunnel config: %v", err)
	}
	if cfg.KnownHostsFile != knownHostsPath {
		t.Fatalf("expected resolved known hosts path %q, got %q", knownHostsPath, cfg.KnownHostsFile)
	}
}
