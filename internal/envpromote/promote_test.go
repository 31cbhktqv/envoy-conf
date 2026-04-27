package envpromote

import (
	"testing"
)

func TestPromote_Ready(t *testing.T) {
	from := Stage{Name: "staging", Env: map[string]string{"A": "1", "B": "2"}}
	to := Stage{Name: "production", Env: map[string]string{"A": "1", "B": "2"}}
	r := Promote(from, to)
	if !r.Ready {
		t.Errorf("expected ready, got not ready")
	}
	if len(r.MissingKeys) != 0 {
		t.Errorf("expected no missing keys, got %v", r.MissingKeys)
	}
}

func TestPromote_MissingKeys(t *testing.T) {
	from := Stage{Name: "staging", Env: map[string]string{"A": "1", "B": "2", "C": "3"}}
	to := Stage{Name: "production", Env: map[string]string{"A": "1"}}
	r := Promote(from, to)
	if r.Ready {
		t.Error("expected not ready due to missing keys")
	}
	if len(r.MissingKeys) != 2 {
		t.Errorf("expected 2 missing keys, got %d: %v", len(r.MissingKeys), r.MissingKeys)
	}
}

func TestPromote_NewKeys(t *testing.T) {
	from := Stage{Name: "staging", Env: map[string]string{"A": "1"}}
	to := Stage{Name: "production", Env: map[string]string{"A": "1", "B": "extra"}}
	r := Promote(from, to)
	if !r.Ready {
		t.Error("expected ready; new keys in target should not block promotion")
	}
	if len(r.NewKeys) != 1 || r.NewKeys[0] != "B" {
		t.Errorf("expected NewKeys=[B], got %v", r.NewKeys)
	}
}

func TestPromote_ChangedKeys(t *testing.T) {
	from := Stage{Name: "staging", Env: map[string]string{"A": "1", "B": "old"}}
	to := Stage{Name: "production", Env: map[string]string{"A": "1", "B": "new"}}
	r := Promote(from, to)
	if !r.Ready {
		t.Error("expected ready; changed values should not block promotion")
	}
	if len(r.ChangedKeys) != 1 || r.ChangedKeys[0] != "B" {
		t.Errorf("expected ChangedKeys=[B], got %v", r.ChangedKeys)
	}
}

func TestSummary_Ready(t *testing.T) {
	r := PromoteResult{From: "staging", To: "production", Ready: true, ChangedKeys: []string{"X"}, NewKeys: []string{"Y"}}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_NotReady(t *testing.T) {
	r := PromoteResult{From: "staging", To: "production", Ready: false, MissingKeys: []string{"DB_URL"}}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
