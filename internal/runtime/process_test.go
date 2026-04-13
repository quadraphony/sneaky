package runtime

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestStartProcessAndStop(t *testing.T) {
	binary := writeProcessScript(t)
	cmd := exec.Command(binary)

	handle, err := StartProcess(cmd, nil)
	if err != nil {
		t.Fatalf("start process: %v", err)
	}
	if handle.State() != StateRunning {
		t.Fatalf("expected running state, got %q", handle.State())
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := handle.Stop(stopCtx); err != nil {
		t.Fatalf("stop process: %v", err)
	}
	if handle.State() != StateStopped {
		t.Fatalf("expected stopped state, got %q", handle.State())
	}
}

func writeProcessScript(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "worker.sh")
	script := `#!/usr/bin/env bash
set -euo pipefail
trap 'exit 0' TERM INT
while true; do
  sleep 1
done
`
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write process script: %v", err)
	}
	return path
}
