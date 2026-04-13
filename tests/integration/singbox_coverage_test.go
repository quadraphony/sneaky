package integration

import (
	"context"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"sneaky-core/internal/tools"
)

type protocolCase struct {
	name       string
	serverPath string
	clientPath string
}

func TestSingboxCoverageFirstBatch(t *testing.T) {
	singboxBin, err := tools.ResolveSingbox()
	if err != nil {
		t.Skip("sing-box binary not available")
	}

	repoRoot := locateRepoRoot(t)
	cases := []protocolCase{
		{name: "vless", serverPath: "testdata/singbox/vless/server.json", clientPath: "testdata/singbox/vless/client.json"},
		{name: "vmess", serverPath: "testdata/singbox/vmess/server.json", clientPath: "testdata/singbox/vmess/client.json"},
		{name: "trojan", serverPath: "testdata/singbox/trojan/server.json", clientPath: "testdata/singbox/trojan/client.json"},
		{name: "shadowsocks", serverPath: "testdata/singbox/shadowsocks/server.json", clientPath: "testdata/singbox/shadowsocks/client.json"},
		{name: "hysteria2", serverPath: "testdata/singbox/hysteria2/server.json", clientPath: "testdata/singbox/hysteria2/client.json"},
		{name: "tuic", serverPath: "testdata/singbox/tuic/server.json", clientPath: "testdata/singbox/tuic/client.json"},
		{name: "hysteria", serverPath: "testdata/singbox/hysteria/server.json", clientPath: "testdata/singbox/hysteria/client.json"},
		{name: "http-connect", serverPath: "testdata/singbox/http/server.json", clientPath: "testdata/singbox/http/client.json"},
		{name: "shadowtls", serverPath: "testdata/singbox/shadowtls/server.json", clientPath: "testdata/singbox/shadowtls/client.json"},
		{name: "anytls", serverPath: "testdata/singbox/anytls/server.json", clientPath: "testdata/singbox/anytls/client.json"},
		{name: "naive", serverPath: "testdata/singbox/naive/server.json", clientPath: "testdata/singbox/naive/client.json"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			serverConfig := filepath.Join(repoRoot, tc.serverPath)
			clientConfig := filepath.Join(repoRoot, tc.clientPath)

			checkConfig(t, repoRoot, singboxBin, serverConfig)
			checkConfig(t, repoRoot, singboxBin, clientConfig)

			server := exec.Command(singboxBin, "run", "-c", serverConfig, "--disable-color")
			server.Dir = repoRoot
			if err := server.Start(); err != nil {
				t.Fatalf("start server %s: %v", tc.name, err)
			}
			defer func() {
				if server.Process != nil {
					_ = server.Process.Kill()
					_, _ = server.Process.Wait()
				}
			}()

			time.Sleep(400 * time.Millisecond)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			probe := exec.CommandContext(ctx, "go", "run", "./cmd/sneakycli", "probe", clientConfig, "https://example.com")
			probe.Dir = repoRoot
			output, err := probe.CombinedOutput()
			if err != nil {
				t.Fatalf("probe %s failed: %v\n%s", tc.name, err, output)
			}

			if want := "probe ok adapter=singbox"; !strings.Contains(string(output), want) {
				t.Fatalf("unexpected probe output for %s: %s", tc.name, output)
			}
		})
	}
}

func TestSingboxUtilityFixturesValidate(t *testing.T) {
	singboxBin, err := tools.ResolveSingbox()
	if err != nil {
		t.Skip("sing-box binary not available")
	}

	repoRoot := locateRepoRoot(t)
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/utilities/dns.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/utilities/routing.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/shadowtls/server.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/shadowtls/client.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/anytls/server.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/anytls/client.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/naive/server.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/naive/client.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/ssh/client.json"))
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/tor/client.json"))
}

func TestSingboxSSHProbe(t *testing.T) {
	singboxBin, err := tools.ResolveSingbox()
	if err != nil {
		t.Skip("sing-box binary not available")
	}

	repoRoot := locateRepoRoot(t)
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/ssh/client.json"))

	sshd := startLocalSSHD(t, repoRoot)
	defer func() {
		if sshd.Process != nil {
			_ = sshd.Process.Kill()
			_, _ = sshd.Process.Wait()
		}
	}()

	output := runProbe(t, repoRoot, "testdata/singbox/ssh/client.json", 30*time.Second)
	if want := "probe ok adapter=singbox"; !strings.Contains(output, want) {
		t.Fatalf("unexpected SSH probe output: %s", output)
	}
}

func TestSingboxTorProbe(t *testing.T) {
	singboxBin, err := tools.ResolveSingbox()
	if err != nil {
		t.Skip("sing-box binary not available")
	}

	repoRoot := locateRepoRoot(t)
	checkConfig(t, repoRoot, singboxBin, filepath.Join(repoRoot, "testdata/singbox/tor/client.json"))

	output := runProbe(t, repoRoot, "testdata/singbox/tor/client.json", 60*time.Second)
	if want := "probe ok adapter=singbox"; !strings.Contains(output, want) {
		t.Fatalf("unexpected Tor probe output: %s", output)
	}
}

func checkConfig(t *testing.T, repoRoot, singboxBin, path string) {
	t.Helper()

	cmd := exec.Command(singboxBin, "check", "-c", path, "--disable-color")
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check config %s failed: %v\n%s", path, err, output)
	}
}

func locateRepoRoot(t *testing.T) string {
	t.Helper()

	path, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return path
}

func runProbe(t *testing.T, repoRoot, clientConfig string, timeout time.Duration) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	probe := exec.CommandContext(ctx, "go", "run", "./cmd/sneakycli", "probe", filepath.Join(repoRoot, clientConfig), "https://example.com")
	probe.Dir = repoRoot
	output, err := probe.CombinedOutput()
	if err != nil {
		t.Fatalf("probe %s failed: %v\n%s", clientConfig, err, output)
	}
	return string(output)
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
	if err := os.WriteFile(authorizedKeys, []byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINdWoBKqncxk5tvvso20L/yyvPtBEU1xEfcdsu/jUmJL leeroy@Leeroy\n"), 0o600); err != nil {
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

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:22322", 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return cmd
		}
		time.Sleep(100 * time.Millisecond)
	}

	if cmd.Process != nil {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}
	t.Fatalf("local sshd was not ready in time")
	return nil
}
