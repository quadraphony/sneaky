package ssh

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/proxy"
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

func TestCapabilitiesDeclareSSHFamilyBehavior(t *testing.T) {
	caps := New("ssh").Capabilities()
	if !caps.SupportsLocalSOCKS {
		t.Fatal("expected SSH adapter to support local SOCKS forwarding")
	}
	if !caps.HonorsKnownHostsFile {
		t.Fatal("expected SSH adapter to honor known_hosts_file")
	}
	if !caps.HonorsHostKeyChecking {
		t.Fatal("expected SSH adapter to honor strict host key checking")
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

func TestStartAndProbeThroughLocalSSHD(t *testing.T) {
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		t.Skip("ssh binary not available")
	}

	repoRoot := locateRepoRoot(t)
	sshd := startLocalSSHD(t, repoRoot)
	defer func() {
		if sshd.Process != nil {
			_ = sshd.Process.Kill()
			_, _ = sshd.Process.Wait()
		}
	}()

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("current user: %v", err)
	}

	knownHostsPath := filepath.Join(t.TempDir(), "known_hosts")
	hostKeyLine, err := os.ReadFile(filepath.Join(repoRoot, "testdata/ssh/host_ed25519.pub"))
	if err != nil {
		t.Fatalf("read host public key: %v", err)
	}
	fields := strings.Fields(string(hostKeyLine))
	if len(fields) < 2 {
		t.Fatalf("unexpected host public key format: %q", string(hostKeyLine))
	}
	knownHostsEntry := "[127.0.0.1]:22322 " + fields[0] + " " + fields[1] + "\n"
	if err := os.WriteFile(knownHostsPath, []byte(knownHostsEntry), 0o600); err != nil {
		t.Fatalf("write known hosts: %v", err)
	}

	socksPort := reserveLocalPort(t)
	cfg := `{
	  "ssh_tunnel": {
	    "host": "127.0.0.1",
	    "user": "` + currentUser.Username + `",
	    "port": 22322,
	    "local_socks_port": ` + strconv.Itoa(socksPort) + `,
	    "identity_file": "` + filepath.Join(repoRoot, "testdata/ssh/client_ed25519") + `",
	    "known_hosts_file": "` + knownHostsPath + `",
	    "strict_host_key_checking": "yes",
	    "connect_timeout_seconds": 5
	  }
	}`

	a := New(sshBinary)
	if err := a.ValidateConfig(adapter.StartRequest{RawConfig: []byte(cfg)}); err != nil {
		t.Fatalf("validate config: %v", err)
	}

	handle, err := a.Start(context.Background(), adapter.StartRequest{RawConfig: []byte(cfg)})
	if err != nil {
		t.Fatalf("start adapter: %v", err)
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = handle.Stop(stopCtx)
	}()

	waitForPort(t, socksPort, 5*time.Second)
	statusCode, err := probeViaSOCKS(socksPort, "https://example.com")
	if err != nil {
		t.Fatalf("probe through ssh socks tunnel: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200 through ssh socks tunnel, got %d", statusCode)
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

func TestValidateConfigPassesKnownHostsAndHostKeyOptions(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "args.log")
	knownHostsPath := filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(knownHostsPath, []byte("example"), 0o600); err != nil {
		t.Fatalf("write known hosts: %v", err)
	}

	binary := writeLoggingSSHBinary(t, logPath)
	a := New(binary)

	cfg := `{
	  "ssh_tunnel": {
	    "host": "example.com",
	    "user": "demo",
	    "local_socks_port": 1080,
	    "known_hosts_file": "` + knownHostsPath + `",
	    "strict_host_key_checking": "yes"
	  }
	}`
	if err := a.ValidateConfig(adapter.StartRequest{RawConfig: []byte(cfg)}); err != nil {
		t.Fatalf("validate config: %v", err)
	}

	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read logged args: %v", err)
	}
	args := string(logged)
	if !strings.Contains(args, "UserKnownHostsFile="+knownHostsPath) {
		t.Fatalf("expected known hosts option in args, got %q", args)
	}
	if !strings.Contains(args, "StrictHostKeyChecking=yes") {
		t.Fatalf("expected strict host key checking option in args, got %q", args)
	}
}

func TestStartPassesKnownHostsAndKeepaliveOptions(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "args.log")
	knownHostsPath := filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(knownHostsPath, []byte("example"), 0o600); err != nil {
		t.Fatalf("write known hosts: %v", err)
	}

	binary := writeLoggingSSHBinary(t, logPath)
	a := New(binary)

	cfg := `{
	  "ssh_tunnel": {
	    "host": "example.com",
	    "user": "demo",
	    "local_socks_port": 1080,
	    "known_hosts_file": "` + knownHostsPath + `",
	    "strict_host_key_checking": "yes",
	    "connect_timeout_seconds": 5,
	    "server_alive_interval_seconds": 7,
	    "server_alive_count_max": 2
	  }
	}`

	handle, err := a.Start(context.Background(), adapter.StartRequest{RawConfig: []byte(cfg)})
	if err != nil {
		t.Fatalf("start adapter: %v", err)
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := handle.Stop(stopCtx); err != nil {
		t.Fatalf("stop adapter: %v", err)
	}

	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read logged args: %v", err)
	}
	args := string(logged)
	for _, want := range []string{
		"UserKnownHostsFile=" + knownHostsPath,
		"StrictHostKeyChecking=yes",
		"ConnectTimeout=5",
		"ServerAliveInterval=7",
		"ServerAliveCountMax=2",
		"-D",
		"127.0.0.1:1080",
	} {
		if !strings.Contains(args, want) {
			t.Fatalf("expected %q in args, got %q", want, args)
		}
	}
}

func writeLoggingSSHBinary(t *testing.T, logPath string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "ssh")
	script := `#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" > "` + logPath + `"
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

func locateRepoRoot(t *testing.T) string {
	t.Helper()

	path, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return path
}

func startLocalSSHD(t *testing.T, repoRoot string) *exec.Cmd {
	t.Helper()

	sshdBin := filepath.Join(repoRoot, ".tools/openssh/usr/sbin/sshd")
	libwrapDir := filepath.Join(repoRoot, ".tools/libwrap/usr/lib/x86_64-linux-gnu")
	if _, err := os.Stat(sshdBin); err != nil {
		t.Skip("local sshd binary not available")
	}
	if _, err := os.Stat(libwrapDir); err != nil {
		t.Skip("local libwrap runtime not available")
	}

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("current user: %v", err)
	}

	tmpDir := t.TempDir()
	authorizedKeys := filepath.Join(tmpDir, "authorized_keys")
	pubKeyPath := filepath.Join(repoRoot, "testdata/ssh/client_ed25519.pub")
	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		t.Fatalf("read client public key: %v", err)
	}
	if err := os.WriteFile(authorizedKeys, pubKey, 0o600); err != nil {
		t.Fatalf("write authorized_keys: %v", err)
	}

	configPath := filepath.Join(tmpDir, "sshd_config")
	config := strings.Join([]string{
		"Port 22322",
		"ListenAddress 127.0.0.1",
		"HostKey " + filepath.Join(repoRoot, "testdata/ssh/host_ed25519"),
		"AuthorizedKeysFile " + authorizedKeys,
		"PubkeyAuthentication yes",
		"PasswordAuthentication no",
		"KbdInteractiveAuthentication no",
		"ChallengeResponseAuthentication no",
		"UsePAM no",
		"PermitRootLogin no",
		"AllowUsers " + currentUser.Username,
		"PidFile " + filepath.Join(tmpDir, "sshd.pid"),
		"LogLevel VERBOSE",
		"StrictModes no",
		"Subsystem sftp " + filepath.Join(repoRoot, ".tools/openssh/usr/lib/openssh/sftp-server"),
	}, "\n") + "\n"
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write sshd config: %v", err)
	}

	cmd := exec.Command(sshdBin, "-D", "-e", "-f", configPath)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "LD_LIBRARY_PATH="+libwrapDir)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start local sshd: %v", err)
	}

	waitForPort(t, 22322, 5*time.Second)
	return cmd
}

func reserveLocalPort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve local port: %v", err)
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener addr type: %T", listener.Addr())
	}
	return addr.Port
}

func waitForPort(t *testing.T, port int, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("port %s was not ready within %s", address, timeout)
}

func probeViaSOCKS(port int, targetURL string) (int, error) {
	proxyAddr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	socksDialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return 0, err
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		type contextDialer interface {
			DialContext(context.Context, string, string) (net.Conn, error)
		}
		if d, ok := socksDialer.(contextDialer); ok {
			return d.DialContext(ctx, network, addr)
		}
		return socksDialer.Dial(network, addr)
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			DialContext:         dialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if _, err := io.ReadAll(resp.Body); err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
