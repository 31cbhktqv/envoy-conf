package auditor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadLog(t *testing.T) {
	a := New()
	a.Record(EventDiff, []string{"staging", "prod"}, map[string]string{"keys": "3"}, true, "diff ok")
	a.Record(EventValidate, []string{"dev"}, nil, false, "pattern mismatch")

	tmp := filepath.Join(t.TempDir(), "audit", "log.json")
	if err := a.SaveLog(tmp); err != nil {
		t.Fatalf("SaveLog: %v", err)
	}

	loaded, err := LoadLog(tmp)
	if err != nil {
		t.Fatalf("LoadLog: %v", err)
	}

	if len(loaded.Entries()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries()))
	}

	e0 := loaded.Entries()[0]
	if e0.Event != EventDiff {
		t.Errorf("expected event %q, got %q", EventDiff, e0.Event)
	}
	if e0.Meta["keys"] != "3" {
		t.Errorf("expected meta keys=3, got %q", e0.Meta["keys"])
	}

	e1 := loaded.Entries()[1]
	if e1.Success {
		t.Error("expected success=false for second entry")
	}
}

func TestLoadLog_MissingFile(t *testing.T) {
	_, err := LoadLog("/nonexistent/path/audit.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadLog_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(tmp, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadLog(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveLog_CreatesDirectory(t *testing.T) {
	a := New()
	a.Record(EventSnapshot, nil, nil, true, "")

	dir := filepath.Join(t.TempDir(), "nested", "dir")
	path := filepath.Join(dir, "audit.json")

	if err := a.SaveLog(path); err != nil {
		t.Fatalf("SaveLog: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
