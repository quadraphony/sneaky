package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	envSingboxBin = "SNEAKY_SINGBOX_BIN"
)

func ResolveSingbox() (string, error) {
	if path := os.Getenv(envSingboxBin); path != "" {
		return filepath.Abs(path)
	}

	candidates := []string{
		".tools/sing-box/current/sing-box",
		".tools/sing-box/1.13.7/sing-box-1.13.7-linux-amd64/sing-box",
	}

	wd, err := os.Getwd()
	if err == nil {
		for _, base := range walkUp(wd) {
			for _, candidate := range candidates {
				abs := filepath.Join(base, candidate)
				if info, err := os.Stat(abs); err == nil && !info.IsDir() && info.Mode()&0o111 != 0 {
					return abs, nil
				}
			}
		}
	}

	path, err := exec.LookPath("sing-box")
	if err != nil {
		return "", fmt.Errorf("resolve sing-box binary: %w", err)
	}
	return path, nil
}

func walkUp(start string) []string {
	paths := []string{start}
	for {
		parent := filepath.Dir(start)
		if parent == start {
			return paths
		}
		paths = append(paths, parent)
		start = parent
	}
}
