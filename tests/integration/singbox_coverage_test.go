package integration

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type protocolCase struct {
	name       string
	serverPath string
	clientPath string
}

func TestSingboxCoverageFirstBatch(t *testing.T) {
	if _, err := exec.LookPath("sing-box"); err != nil {
		t.Skip("sing-box binary not available")
	}

	repoRoot := locateRepoRoot(t)
	cases := []protocolCase{
		{name: "vless", serverPath: "testdata/singbox/vless/server.json", clientPath: "testdata/singbox/vless/client.json"},
		{name: "vmess", serverPath: "testdata/singbox/vmess/server.json", clientPath: "testdata/singbox/vmess/client.json"},
		{name: "trojan", serverPath: "testdata/singbox/trojan/server.json", clientPath: "testdata/singbox/trojan/client.json"},
		{name: "shadowsocks", serverPath: "testdata/singbox/shadowsocks/server.json", clientPath: "testdata/singbox/shadowsocks/client.json"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			serverConfig := filepath.Join(repoRoot, tc.serverPath)
			clientConfig := filepath.Join(repoRoot, tc.clientPath)

			checkConfig(t, repoRoot, serverConfig)
			checkConfig(t, repoRoot, clientConfig)

			server := exec.Command("sing-box", "run", "-c", serverConfig, "--disable-color")
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

func checkConfig(t *testing.T, repoRoot, path string) {
	t.Helper()

	cmd := exec.Command("sing-box", "check", "-c", path, "--disable-color")
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
