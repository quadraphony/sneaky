package singbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/config"
)

func (a *Adapter) prepareConfig(req adapter.StartRequest) (string, func(), error) {
	if len(req.RawConfig) > 0 {
		input, err := config.Parse(req.RawConfig, "inline")
		if err != nil {
			return "", noopCleanup, err
		}
		metadata, err := config.DetectAndValidate(input)
		if err != nil {
			return "", noopCleanup, err
		}
		if metadata.AdapterID != config.AdapterSingbox {
			return "", noopCleanup, fmt.Errorf("sing-box adapter received unsupported adapter id %q", metadata.AdapterID)
		}

		tmpFile, err := os.CreateTemp("", "sneaky-singbox-*.json")
		if err != nil {
			return "", noopCleanup, fmt.Errorf("create temp config: %w", err)
		}
		if _, err := tmpFile.Write(bytes.TrimSpace(req.RawConfig)); err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return "", noopCleanup, fmt.Errorf("write temp config: %w", err)
		}
		if err := tmpFile.Close(); err != nil {
			os.Remove(tmpFile.Name())
			return "", noopCleanup, fmt.Errorf("close temp config: %w", err)
		}

		return tmpFile.Name(), func() {
			_ = os.Remove(tmpFile.Name())
		}, nil
	}

	if req.ConfigPath == "" {
		return "", noopCleanup, fmt.Errorf("sing-box adapter requires config bytes or a config path")
	}

	absPath, err := filepath.Abs(req.ConfigPath)
	if err != nil {
		return "", noopCleanup, fmt.Errorf("resolve config path: %w", err)
	}

	input, err := config.LoadFile(absPath)
	if err != nil {
		return "", noopCleanup, err
	}
	metadata, err := config.DetectAndValidate(input)
	if err != nil {
		return "", noopCleanup, err
	}
	if metadata.AdapterID != config.AdapterSingbox {
		return "", noopCleanup, fmt.Errorf("config at %s is not a sing-box config", absPath)
	}

	return absPath, noopCleanup, nil
}

func (a *Adapter) checkConfig(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, a.binary(), "check", "-c", path, "--disable-color")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			return fmt.Errorf("sing-box check failed: %w", err)
		}
		return fmt.Errorf("sing-box check failed: %s", bytes.TrimSpace(output))
	}
	return nil
}

func (a *Adapter) binary() string {
	if a.binaryPath != "" {
		return a.binaryPath
	}
	return "sing-box"
}

func noopCleanup() {}
