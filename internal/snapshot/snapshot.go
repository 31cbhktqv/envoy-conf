package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of environment variables at a point in time.
type Snapshot struct {
	Target    string            `json:"target"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot for the given target and env map.
func New(target string, env map[string]string) *Snapshot {
	return &Snapshot{
		Target:    target,
		Timestamp: time.Now().UTC(),
		Env:       env,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(s *Snapshot, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	if s.Target == "" {
		return nil, fmt.Errorf("snapshot: missing target field")
	}
	if s.Env == nil {
		s.Env = make(map[string]string)
	}
	return &s, nil
}
