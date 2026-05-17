package ssh

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/config"
	"sneaky-core/internal/runtime"
)

type Adapter struct {
	binaryPath string
}

func New(binaryPath string) *Adapter {
	return &Adapter{binaryPath: binaryPath}
}

func (a *Adapter) Identity() string {
	return config.AdapterSSH
}

func (a *Adapter) Capabilities() adapter.Capabilities {
	return adapter.Capabilities{
		SupportsLocalSOCKS:    true,
		HonorsKnownHostsFile:  true,
		HonorsHostKeyChecking: true,
	}
}

func (a *Adapter) ValidateConfig(req adapter.StartRequest) error {
	bin, err := a.resolveBinary()
	if err != nil {
		return err
	}

	tunnel, err := loadTunnelConfig(req)
	if err != nil {
		return err
	}

	checkArgs := []string{
		"-G",
		"-p", strconv.Itoa(tunnel.Port),
		"-o", "StrictHostKeyChecking=" + tunnel.StrictHostKeyChecking,
		"-o", "ConnectTimeout=" + strconv.Itoa(tunnel.ConnectTimeoutSeconds),
		"-o", "ServerAliveInterval=" + strconv.Itoa(tunnel.ServerAliveInterval),
		"-o", "ServerAliveCountMax=" + strconv.Itoa(tunnel.ServerAliveCountMax),
	}
	if tunnel.IdentityFile != "" {
		checkArgs = append(checkArgs, "-i", tunnel.IdentityFile)
	}
	if tunnel.KnownHostsFile != "" {
		checkArgs = append(checkArgs, "-o", "UserKnownHostsFile="+tunnel.KnownHostsFile)
	}
	checkArgs = append(checkArgs, tunnel.User+"@"+tunnel.Host)

	cmd := exec.Command(bin, checkArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ssh config check failed: %s", output)
	}
	return nil
}

func (a *Adapter) Start(ctx context.Context, req adapter.StartRequest) (runtime.Handle, error) {
	bin, err := a.resolveBinary()
	if err != nil {
		return nil, err
	}

	tunnel, err := loadTunnelConfig(req)
	if err != nil {
		return nil, err
	}

	localAddr := "127.0.0.1:" + strconv.Itoa(tunnel.LocalSOCKSPort)

	args := []string{
		"-N",
		"-D", localAddr,
		"-p", strconv.Itoa(tunnel.Port),
		"-o", "ExitOnForwardFailure=yes",
		"-o", "StrictHostKeyChecking=" + tunnel.StrictHostKeyChecking,
		"-o", "ConnectTimeout=" + strconv.Itoa(tunnel.ConnectTimeoutSeconds),
		"-o", "ServerAliveInterval=" + strconv.Itoa(tunnel.ServerAliveInterval),
		"-o", "ServerAliveCountMax=" + strconv.Itoa(tunnel.ServerAliveCountMax),
		tunnel.User + "@" + tunnel.Host,
	}
	if tunnel.IdentityFile != "" {
		args = append(args[:len(args)-1], append([]string{"-i", tunnel.IdentityFile}, args[len(args)-1])...)
	}
	if tunnel.KnownHostsFile != "" {
		args = append(args[:len(args)-1], append([]string{"-o", "UserKnownHostsFile=" + tunnel.KnownHostsFile}, args[len(args)-1])...)
	}

	ready := func() bool {
		conn, err := net.DialTimeout("tcp", localAddr, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}

	cmd := exec.CommandContext(ctx, bin, args...)
	handle, err := runtime.StartProcessWithReadiness(cmd, nil, ready)
	if err != nil {
		return nil, err
	}
	return handle, nil
}

func (a *Adapter) resolveBinary() (string, error) {
	if a.binaryPath != "" {
		if _, err := os.Stat(a.binaryPath); err != nil {
			return "", fmt.Errorf("ssh binary not found at %s: %w", a.binaryPath, err)
		}
		return a.binaryPath, nil
	}
	path, err := exec.LookPath("ssh")
	if err != nil {
		return "", fmt.Errorf("failed to resolve ssh binary: %w", err)
	}
	return path, nil
}
