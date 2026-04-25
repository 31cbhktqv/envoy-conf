package envcompare_test

import (
	"testing"

	"github.com/yourorg/envoy-conf/internal/envcompare"
)

func TestCompare_AllMatch(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := envcompare.Compare(a, b, envcompare.CompareOptions{CaseSensitiveValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MatchedCount != 2 {
		t.Errorf("expected 2 matched, got %d", res.MatchedCount)
	}
	if len(res.Mismatched) != 0 || len(res.MissingInA) != 0 || len(res.MissingInB) != 0 {
		t.Errorf("expected no differences")
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}
	res, err := envcompare.Compare(a, b, envcompare.CompareOptions{CaseSensitiveValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.MissingInB) != 1 || res.MissingInB[0] != "ONLY_A" {
		t.Errorf("expected ONLY_A missing in B, got %v", res.MissingInB)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "val"}
	res, err := envcompare.Compare(a, b, envcompare.CompareOptions{CaseSensitiveValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.MissingInA) != 1 || res.MissingInA[0] != "ONLY_B" {
		t.Errorf("expected ONLY_B missing in A, got %v", res.MissingInA)
	}
}

func TestCompare_Mismatch(t *testing.T) {
	a := map[string]string{"FOO": "alpha"}
	b := map[string]string{"FOO": "beta"}
	res, err := envcompare.Compare(a, b, envcompare.CompareOptions{CaseSensitiveValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pair, ok := res.Mismatched["FOO"]
	if !ok {
		t.Fatal("expected FOO in mismatched")
	}
	if pair[0] != "alpha" || pair[1] != "beta" {
		t.Errorf("unexpected mismatch values: %v", pair)
	}
}

func TestCompare_IgnoreKeys(t *testing.T) {
	a := map[string]string{"FOO": "x", "SECRET": "a"}
	b := map[string]string{"FOO": "x", "SECRET": "b"}
	opts := envcompare.CompareOptions{IgnoreKeys: []string{"SECRET"}, CaseSensitiveValues: true}
	res, err := envcompare.Compare(a, b, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Mismatched) != 0 {
		t.Errorf("expected SECRET to be ignored, got mismatches: %v", res.Mismatched)
	}
}

func TestCompare_IgnorePattern(t *testing.T) {
	a := map[string]string{"DB_HOST": "localhost", "DB_PASS": "secret1"}
	b := map[string]string{"DB_HOST": "localhost", "DB_PASS": "secret2"}
	opts := envcompare.CompareOptions{IgnorePatterns: []string{"^DB_PASS"}, CaseSensitiveValues: true}
	res, err := envcompare.Compare(a, b, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Mismatched) != 0 {
		t.Errorf("expected DB_PASS to be ignored")
	}
}

func TestCompare_InvalidPattern(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar"}
	opts := envcompare.CompareOptions{IgnorePatterns: []string{"[invalid"}}
	_, err := envcompare.Compare(a, b, opts)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestCompare_CaseInsensitiveValues(t *testing.T) {
	a := map[string]string{"MODE": "Production"}
	b := map[string]string{"MODE": "production"}
	res, err := envcompare.Compare(a, b, envcompare.CompareOptions{CaseSensitiveValues: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Mismatched) != 0 {
		t.Errorf("expected case-insensitive match for MODE")
	}
	if res.MatchedCount != 1 {
		t.Errorf("expected MatchedCount 1, got %d", res.MatchedCount)
	}
}
