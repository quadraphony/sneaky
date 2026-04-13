package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type stateFile struct {
	PID        int       `json:"pid"`
	AdapterID  string    `json:"adapter_id"`
	ConfigPath string    `json:"config_path"`
	StartedAt  time.Time `json:"started_at"`
}

func statePath() string {
	return filepath.Join(os.TempDir(), "sneaky-cli-state.json")
}

func writeStateFile(state stateFile) error {
	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state file: %w", err)
	}

	if err := os.WriteFile(statePath(), raw, 0o600); err != nil {
		return fmt.Errorf("write state file: %w", err)
	}
	return nil
}

func readStateFile() (stateFile, error) {
	raw, err := os.ReadFile(statePath())
	if err != nil {
		return stateFile{}, err
	}

	var state stateFile
	if err := json.Unmarshal(raw, &state); err != nil {
		return stateFile{}, fmt.Errorf("decode state file: %w", err)
	}
	return state, nil
}

func removeStateFile() error {
	err := os.Remove(statePath())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove state file: %w", err)
	}
	return nil
}
