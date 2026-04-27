package envrollout

import (
	"testing"
)

func TestPlan_SingleTransition_Ready(t *testing.T) {
	stages := []Stage{
		{Name: "staging", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "prod", Env: map[string]string{"A": "1", "B": "2"}},
	}
	results := Plan(stages)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Ready {
		t.Errorf("expected ready, got not ready")
	}
}

func TestPlan_MissingKey_NotReady(t *testing.T) {
	stages := []Stage{
		{Name: "staging", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "prod", Env: map[string]string{"A": "1"}},
	}
	results := Plan(stages)
	if results[0].Ready {
		t.Error("expected not ready due to missing key")
	}
	if len(results[0].MissingKeys) != 1 || results[0].MissingKeys[0] != "B" {
		t.Errorf("unexpected missing keys: %v", results[0].MissingKeys)
	}
}

func TestPlan_NewKey_StillReady(t *testing.T) {
	stages := []Stage{
		{Name: "staging", Env: map[string]string{"A": "1"}},
		{Name: "prod", Env: map[string]string{"A": "1", "B": "2"}},
	}
	results := Plan(stages)
	if !results[0].Ready {
		t.Error("expected ready; new keys should not block")
	}
	if len(results[0].NewKeys) != 1 || results[0].NewKeys[0] != "B" {
		t.Errorf("unexpected new keys: %v", results[0].NewKeys)
	}
}

func TestPlan_ChangedKey_StillReady(t *testing.T) {
	stages := []Stage{
		{Name: "staging", Env: map[string]string{"A": "old"}},
		{Name: "prod", Env: map[string]string{"A": "new"}},
	}
	results := Plan(stages)
	if !results[0].Ready {
		t.Error("expected ready; changed keys should not block")
	}
	if len(results[0].ChangedKeys) != 1 || results[0].ChangedKeys[0] != "A" {
		t.Errorf("unexpected changed keys: %v", results[0].ChangedKeys)
	}
}

func TestPlan_MultiStage(t *testing.T) {
	stages := []Stage{
		{Name: "dev", Env: map[string]string{"X": "1"}},
		{Name: "staging", Env: map[string]string{"X": "1"}},
		{Name: "prod", Env: map[string]string{"X": "1"}},
	}
	results := Plan(stages)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestHasBlocker(t *testing.T) {
	results := []Result{
		{Ready: true},
		{Ready: false},
	}
	if !HasBlocker(results) {
		t.Error("expected HasBlocker to return true")
	}
	results[1].Ready = true
	if HasBlocker(results) {
		t.Error("expected HasBlocker to return false")
	}
}
