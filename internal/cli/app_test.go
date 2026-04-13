package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testSingboxConfig = `{
  "log": {"level": "warn"},
  "outbounds": [
    {"type": "direct", "tag": "direct"}
  ]
}`

func TestAppValidateAndVersion(t *testing.T) {
	path := writeTestConfig(t)
	defer removeStateFile()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app := New(&stdout, &stderr)

	if code := app.Run([]string{"validate", path}); code != 0 {
		t.Fatalf("validate exit code = %d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "valid adapter=singbox format=json") {
		t.Fatalf("unexpected validate output: %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"version"}); code != 0 {
		t.Fatalf("version exit code = %d", code)
	}
	if strings.TrimSpace(stdout.String()) != version {
		t.Fatalf("unexpected version output: %q", stdout.String())
	}
}

func TestCLIStartStatusStop(t *testing.T) {
	if _, err := exec.LookPath("sing-box"); err != nil {
		t.Skip("sing-box binary not available")
	}

	_ = removeStateFile()
	defer removeStateFile()

	path := writeTestConfig(t)

	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "start", path)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	var startOut bytes.Buffer
	var startErr bytes.Buffer
	cmd.Stdout = &startOut
	cmd.Stderr = &startErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}

	waitForRunningState(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app := New(&stdout, &stderr)
	if code := app.Run([]string{"status"}); code != 0 {
		t.Fatalf("status exit code = %d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "state=running") {
		t.Fatalf("unexpected status output: %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := app.Run([]string{"stop"}); code != 0 {
		t.Fatalf("stop exit code = %d stderr=%s", code, stderr.String())
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for helper process: %v stdout=%s stderr=%s", err, startOut.String(), startErr.String())
	}
	if !strings.Contains(startOut.String(), "running adapter=singbox") {
		t.Fatalf("unexpected start output: %q", startOut.String())
	}
	if !strings.Contains(startOut.String(), "stopped") {
		t.Fatalf("expected stopped output after signal, got %q", startOut.String())
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}

	app := New(os.Stdout, os.Stderr)
	os.Exit(app.Run(args))
}

func TestPrepareSingboxProbeConfig(t *testing.T) {
	path := writeTestConfig(t)

	raw, port, err := prepareSingboxProbeConfig(path)
	if err != nil {
		t.Fatalf("prepare probe config: %v", err)
	}
	if port <= 0 {
		t.Fatalf("expected positive port, got %d", port)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("decode probe config: %v", err)
	}

	inbounds, ok := doc["inbounds"].([]any)
	if !ok || len(inbounds) != 1 {
		t.Fatalf("expected single probe inbound, got %#v", doc["inbounds"])
	}

	route, ok := doc["route"].(map[string]any)
	if !ok || route["final"] == "" {
		t.Fatalf("expected route final in probe config, got %#v", doc["route"])
	}
}

func writeTestConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(testSingboxConfig), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func waitForRunningState(t *testing.T) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		state, err := readStateFile()
		if err == nil && state.PID > 0 && state.AdapterID != "" {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	state, err := readStateFile()
	t.Fatalf("state file was not created in time: state=%+v err=%v", state, err)
}
