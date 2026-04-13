package ssh

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

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

	cmd := exec.Command(a.binary(), checkArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ssh config check failed: %s", output)
	}
	return nil
}

func (a *Adapter) Start(ctx context.Context, req adapter.StartRequest) (runtime.Handle, error) {
	tunnel, err := loadTunnelConfig(req)
	if err != nil {
		return nil, err
	}

	args := []string{
		"-N",
		"-D", "127.0.0.1:" + strconv.Itoa(tunnel.LocalSOCKSPort),
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

	cmd := exec.CommandContext(ctx, a.binary(), args...)
	handle, err := runtime.StartProcess(cmd, nil)
	if err != nil {
		return nil, fmt.Errorf("start ssh process: %w", err)
	}
	return handle, nil
}

func (a *Adapter) binary() string {
	if a.binaryPath != "" {
		return a.binaryPath
	}
	return "ssh"
}
