package singbox

import (
	"context"
	"fmt"
	"os/exec"

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
	return config.AdapterSingbox
}

func (a *Adapter) Capabilities() adapter.Capabilities {
	return adapter.Capabilities{}
}

func (a *Adapter) ValidateConfig(req adapter.StartRequest) error {
	resolvedPath, cleanup, err := a.prepareConfig(req)
	if err != nil {
		return err
	}
	defer cleanup()

	return a.checkConfig(context.Background(), resolvedPath)
}

func (a *Adapter) Start(ctx context.Context, req adapter.StartRequest) (runtime.Handle, error) {
	resolvedPath, cleanup, err := a.prepareConfig(req)
	if err != nil {
		return nil, err
	}

	if err := a.checkConfig(ctx, resolvedPath); err != nil {
		cleanup()
		return nil, err
	}

	cmd := exec.CommandContext(ctx, a.binary(), "run", "-c", resolvedPath, "--disable-color")
	handle, err := runtime.StartProcess(cmd, cleanup)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("start sing-box process: %w", err)
	}

	return handle, nil
}
