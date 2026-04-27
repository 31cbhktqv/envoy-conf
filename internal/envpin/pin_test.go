package envpin

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
		"PORT":     "8080",
	}
}

func TestPin_StoresTargetAndVars(t *testing.T) {
	env := baseEnv()
	p := Pin("prod", env)
	if p.Target != "prod" {
		t.Errorf("expected target prod, got %s", p.Target)
	}
	if p.Variables["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", p.Variables["PORT"])
	}
	if p.PinnedAt.IsZero() {
		t.Error("expected non-zero PinnedAt")
	}
}

func TestPin_IsolatesCopy(t *testing.T) {
	env := baseEnv()
	p := Pin("prod", env)
	env["PORT"] = "9999"
	if p.Variables["PORT"] != "8080" {
		t.Error("pin should not reflect mutations to original map")
	}
}

func TestCompare_NoChanges(t *testing.T) {
	p := Pin("prod", baseEnv())
	entries := Compare(p, baseEnv())
	if len(entries) != 0 {
		t.Errorf("expected no drift, got %d entries", len(entries))
	}
}

func TestCompare_Changed(t *testing.T) {
	p := Pin("prod", baseEnv())
	current := baseEnv()
	current["LOG_LEVEL"] = "debug"
	entries := Compare(p, current)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != "changed" || entries[0].Key != "LOG_LEVEL" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestCompare_Added(t *testing.T) {
	p := Pin("prod", baseEnv())
	current := baseEnv()
	current["NEW_KEY"] = "value"
	entries := Compare(p, current)
	if len(entries) != 1 || entries[0].Status != "added" {
		t.Errorf("expected 1 added entry, got %+v", entries)
	}
}

func TestCompare_Removed(t *testing.T) {
	p := Pin("prod", baseEnv())
	current := baseEnv()
	delete(current, "PORT")
	entries := Compare(p, current)
	if len(entries) != 1 || entries[0].Status != "removed" {
		t.Errorf("expected 1 removed entry, got %+v", entries)
	}
}

func TestSummary_NoDrift(t *testing.T) {
	s := Summary("prod", []DriftEntry{})
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_WithDrift(t *testing.T) {
	entries := []DriftEntry{
		{Key: "A", Status: "added"},
		{Key: "B", Status: "removed"},
		{Key: "C", Status: "changed"},
	}
	s := Summary("staging", entries)
	expected := "[staging] drift detected: +1 added, -1 removed, ~1 changed"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
