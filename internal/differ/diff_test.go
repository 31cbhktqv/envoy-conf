package differ

import (
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := Diff(a, b)

	if result.HasDifferences() {
		t.Error("expected no differences")
	}
	if len(result.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(result.Unchanged))
	}
}

func TestDiff_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}
	result := Diff(a, b)

	if !result.HasDifferences() {
		t.Error("expected differences")
	}
	if v, ok := result.OnlyInA["ONLY_A"]; !ok || v != "val" {
		t.Errorf("expected ONLY_A=val in OnlyInA, got %v", result.OnlyInA)
	}
}

func TestDiff_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "val"}
	result := Diff(a, b)

	if v, ok := result.OnlyInB["ONLY_B"]; !ok || v != "val" {
		t.Errorf("expected ONLY_B=val in OnlyInB, got %v", result.OnlyInB)
	}
}

func TestDiff_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	result := Diff(a, b)

	pair, ok := result.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected [old new], got %v", pair)
	}
}

func TestDiff_SortedKeys(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "2", "M": "3"}
	b := map[string]string{"Z": "x", "A": "x", "M": "x"}
	result := Diff(a, b)
	keys := result.SortedChangedKeys()

	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys [A M Z], got %v", keys)
	}
}
