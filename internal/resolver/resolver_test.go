package resolver

import (
	"os"
	"testing"
)

func TestResolve_SingleSource(t *testing.T) {
	sources := []Source{
		{Name: "base", Vars: map[string]string{"FOO": "bar", "BAZ": "qux"}},
	}
	result, err := Resolve(sources, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestResolve_LaterSourceOverrides(t *testing.T) {
	sources := []Source{
		{Name: "base", Vars: map[string]string{"KEY": "original"}},
		{Name: "override", Vars: map[string]string{"KEY": "overridden"}},
	}
	result, err := Resolve(sources, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "overridden" {
		t.Errorf("expected 'overridden', got %q", result["KEY"])
	}
}

func TestResolve_NoSources_ReturnsError(t *testing.T) {
	_, err := Resolve([]Source{}, ResolveOptions{})
	if err == nil {
		t.Fatal("expected error for empty sources, got nil")
	}
}

func TestResolve_OverrideKeys_FromOS(t *testing.T) {
	os.Setenv("ENVOY_TEST_KEY", "from-os")
	defer os.Unsetenv("ENVOY_TEST_KEY")

	sources := []Source{
		{Name: "file", Vars: map[string]string{"ENVOY_TEST_KEY": "from-file"}},
	}
	result, err := Resolve(sources, ResolveOptions{OverrideKeys: []string{"ENVOY_TEST_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["ENVOY_TEST_KEY"] != "from-os" {
		t.Errorf("expected 'from-os', got %q", result["ENVOY_TEST_KEY"])
	}
}

func TestSourceNames(t *testing.T) {
	sources := []Source{
		{Name: "alpha", Vars: map[string]string{}},
		{Name: "beta", Vars: map[string]string{}},
	}
	names := SourceNames(sources)
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected names: %v", names)
	}
}
