package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := New("production", env)
	if s.Target != "production" {
		t.Errorf("expected target production, got %s", s.Target)
	}
	if len(s.Env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(s.Env))
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := New("staging", map[string]string{"KEY": "value", "PORT": "8080"})
	original.Timestamp = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if err := Save(original, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Target != original.Target {
		t.Errorf("target mismatch: got %s, want %s", loaded.Target, original.Target)
	}
	if loaded.Env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", loaded.Env["KEY"])
	}
	if loaded.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", loaded.Env["PORT"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_MissingTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notarget.json")
	os.WriteFile(path, []byte(`{"env":{"A":"1"}}`), 0644)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing target")
	}
}
