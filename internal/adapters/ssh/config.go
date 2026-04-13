package ssh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/config"
)

type tunnelConfig struct {
	Host                  string `json:"host"`
	User                  string `json:"user"`
	Port                  int    `json:"port"`
	LocalSOCKSPort        int    `json:"local_socks_port"`
	IdentityFile          string `json:"identity_file"`
	StrictHostKeyChecking string `json:"strict_host_key_checking"`
}

func loadTunnelConfig(req adapter.StartRequest) (tunnelConfig, error) {
	var raw []byte
	var source string
	switch {
	case len(req.RawConfig) > 0:
		raw = req.RawConfig
		source = "inline"
	case req.ConfigPath != "":
		input, err := config.LoadFile(req.ConfigPath)
		if err != nil {
			return tunnelConfig{}, err
		}
		raw = input.Raw
		source = req.ConfigPath
	default:
		return tunnelConfig{}, fmt.Errorf("ssh adapter requires config bytes or a config path")
	}

	input, err := config.Parse(raw, source)
	if err != nil {
		return tunnelConfig{}, err
	}
	metadata, err := config.DetectAndValidate(input)
	if err != nil {
		return tunnelConfig{}, err
	}
	if metadata.AdapterID != config.AdapterSSH {
		return tunnelConfig{}, fmt.Errorf("ssh adapter received unsupported adapter id %q", metadata.AdapterID)
	}

	var parsed struct {
		SSHTunnel tunnelConfig `json:"ssh_tunnel"`
	}
	if err := json.Unmarshal(input.Raw, &parsed); err != nil {
		return tunnelConfig{}, fmt.Errorf("decode ssh config: %w", err)
	}

	if parsed.SSHTunnel.Port == 0 {
		parsed.SSHTunnel.Port = 22
	}
	if parsed.SSHTunnel.StrictHostKeyChecking == "" {
		parsed.SSHTunnel.StrictHostKeyChecking = "accept-new"
	}
	if parsed.SSHTunnel.IdentityFile != "" {
		absPath, err := filepath.Abs(parsed.SSHTunnel.IdentityFile)
		if err != nil {
			return tunnelConfig{}, fmt.Errorf("resolve identity_file: %w", err)
		}
		if _, err := os.Stat(absPath); err != nil {
			return tunnelConfig{}, fmt.Errorf("stat identity_file: %w", err)
		}
		parsed.SSHTunnel.IdentityFile = absPath
	}

	return parsed.SSHTunnel, nil
}
