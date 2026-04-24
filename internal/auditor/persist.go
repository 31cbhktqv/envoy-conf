package auditor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveLog writes the audit log to a JSON file at the given path.
func (a *Auditor) SaveLog(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("auditor: create directory: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("auditor: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(a.log); err != nil {
		return fmt.Errorf("auditor: encode log: %w", err)
	}
	return nil
}

// LoadLog reads an audit log from a JSON file and returns an Auditor populated with its entries.
func LoadLog(path string) (*Auditor, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("auditor: open file: %w", err)
	}
	defer f.Close()

	var log Log
	if err := json.NewDecoder(f).Decode(&log); err != nil {
		return nil, fmt.Errorf("auditor: decode log: %w", err)
	}
	return &Auditor{log: log}, nil
}
