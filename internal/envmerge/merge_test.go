package envmerge

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	res, err := Merge(StrategyLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Env["FOO"] != "1" || res.Env["BAR"] != "2" || res.Env["BAZ"] != "3" {
		t.Errorf("unexpected env map: %v", res.Env)
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"FOO": "base"}
	b := map[string]string{"FOO": "override"}

	res, err := Merge(StrategyLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "override" {
		t.Errorf("expected 'override', got %q", res.Env["FOO"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"FOO": "base"}
	b := map[string]string{"FOO": "override"}

	res, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "base" {
		t.Errorf("expected 'base', got %q", res.Env["FOO"])
	}
}

func TestMerge_StrategyStrict_Conflict(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "2"}

	_, err := Merge(StrategyStrict, a, b)
	if err == nil {
		t.Fatal("expected error for strict strategy with conflict, got nil")
	}
}

func TestMerge_StrategyStrict_NoConflict(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}

	res, err := Merge(StrategyStrict, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := Merge(StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
