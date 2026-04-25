package envtemplate_test

import (
	"testing"

	"envoy-conf/internal/envtemplate"
)

// TestExpand_ChainedValues verifies that values resolved from the lookup map
// are NOT themselves expanded (single-pass semantics).
func TestExpand_ChainedValues(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "${C}",
	}
	lookup := map[string]string{
		"B": "resolved-b",
		"C": "resolved-c",
	}
	out, err := envtemplate.Expand(env, lookup, envtemplate.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A should resolve to the literal value of B in lookup, not chain.
	if out["A"] != "resolved-b" {
		t.Errorf("A: expected resolved-b, got %q", out["A"])
	}
	// B should resolve to resolved-c from lookup.
	if out["B"] != "resolved-c" {
		t.Errorf("B: expected resolved-c, got %q", out["B"])
	}
}

func TestExpand_MixedSyntaxInSameValue(t *testing.T) {
	env := map[string]string{
		"URL": "http://$HOST:${PORT}/path",
	}
	lookup := map[string]string{
		"HOST": "example.com",
		"PORT": "8080",
	}
	out, err := envtemplate.Expand(env, lookup, envtemplate.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "http://example.com:8080/path"
	if out["URL"] != want {
		t.Errorf("URL: expected %q, got %q", want, out["URL"])
	}
}

func TestExpand_EmptyLookupNonStrict(t *testing.T) {
	env := map[string]string{
		"KEY": "prefix_${MISSING_VAR}_suffix",
	}
	opts := envtemplate.Options{Strict: false, Fallback: "DEFAULT"}
	out, err := envtemplate.Expand(env, map[string]string{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "prefix_DEFAULT_suffix" {
		t.Errorf("KEY: expected prefix_DEFAULT_suffix, got %q", out["KEY"])
	}
}
