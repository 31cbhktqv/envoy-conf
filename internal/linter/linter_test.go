package linter

import (
	"testing"
)

func TestLint_NoViolations(t *testing.T) {
	l := New(DefaultRules())
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	violations := l.Lint(env)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	l := New(DefaultRules())
	env := map[string]string{
		"app_host": "localhost",
	}
	violations := l.Lint(env)
	if !hasRule(violations, "no-lowercase-key") {
		t.Error("expected no-lowercase-key violation")
	}
}

func TestLint_EmptyValue(t *testing.T) {
	l := New(DefaultRules())
	env := map[string]string{
		"APP_SECRET": "",
	}
	violations := l.Lint(env)
	if !hasRule(violations, "no-empty-value") {
		t.Error("expected no-empty-value violation")
	}
}

func TestLint_WhitespaceInKey(t *testing.T) {
	l := New(DefaultRules())
	env := map[string]string{
		"APP HOST": "value",
	}
	violations := l.Lint(env)
	if !hasRule(violations, "no-whitespace-in-key") {
		t.Error("expected no-whitespace-in-key violation")
	}
}

func TestLint_CustomRule(t *testing.T) {
	rules := []Rule{
		{
			Name:    "no-test-prefix",
			Message: "key must not start with TEST_",
			Check: func(key, _ string) bool {
				return len(key) >= 5 && key[:5] == "TEST_"
			},
		},
	}
	l := New(rules)
	env := map[string]string{
		"TEST_VAR": "value",
		"PROD_VAR": "value",
	}
	violations := l.Lint(env)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "TEST_VAR" {
		t.Errorf("expected violation on TEST_VAR, got %s", violations[0].Key)
	}
}

func hasRule(violations []Violation, name string) bool {
	for _, v := range violations {
		if v.Rule == name {
			return true
		}
	}
	return false
}
