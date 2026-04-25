package envtemplate

import (
	"os"
	"testing"
)

func TestExpand_NoPlaceholders(t *testing.T) {
	env := map[string]string{"FOO": "bar", "NUM": "42"}
	out, err := Expand(env, nil, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["NUM"] != "42" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestExpand_CurlyBraceSyntax(t *testing.T) {
	env := map[string]string{"DSN": "postgres://${DB_HOST}:${DB_PORT}/mydb"}
	lookup := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, err := Expand(env, lookup, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost:5432/mydb" {
		t.Errorf("got %q", out["DSN"])
	}
}

func TestExpand_BareVarSyntax(t *testing.T) {
	env := map[string]string{"GREETING": "Hello $NAME"}
	lookup := map[string]string{"NAME": "World"}
	out, err := Expand(env, lookup, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "Hello World" {
		t.Errorf("got %q", out["GREETING"])
	}
}

func TestExpand_FallsBackToOS(t *testing.T) {
	os.Setenv("_ENVCONF_TEST_VAR", "from-os")
	defer os.Unsetenv("_ENVCONF_TEST_VAR")

	env := map[string]string{"VAL": "${_ENVCONF_TEST_VAR}"}
	out, err := Expand(env, map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "from-os" {
		t.Errorf("got %q", out["VAL"])
	}
}

func TestExpand_UnresolvedNonStrict(t *testing.T) {
	env := map[string]string{"KEY": "${MISSING}"}
	out, err := Expand(env, map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "" {
		t.Errorf("expected empty fallback, got %q", out["KEY"])
	}
}

func TestExpand_UnresolvedStrict(t *testing.T) {
	env := map[string]string{"KEY": "${MISSING}"}
	opts := Options{Strict: true}
	_, err := Expand(env, map[string]string{}, opts)
	if err == nil {
		t.Fatal("expected error for unresolved variable in strict mode")
	}
}

func TestExpand_LookupOverridesOS(t *testing.T) {
	os.Setenv("_ENVCONF_OVER", "os-value")
	defer os.Unsetenv("_ENVCONF_OVER")

	env := map[string]string{"V": "${_ENVCONF_OVER}"}
	lookup := map[string]string{"_ENVCONF_OVER": "lookup-value"}
	out, err := Expand(env, lookup, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["V"] != "lookup-value" {
		t.Errorf("expected lookup-value, got %q", out["V"])
	}
}
