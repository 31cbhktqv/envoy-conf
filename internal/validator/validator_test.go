package validator

import (
	"testing"
)

func TestValidate_AllPass(t *testing.T) {
	env := map[string]string{
		"APP_PORT": "8080",
		"APP_ENV":  "production",
	}
	rules := []Rule{
		{Key: "APP_PORT", Required: true, Pattern: `^\d+$`},
		{Key: "APP_ENV", Required: true},
	}
	violations := Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{
		{Key: "DATABASE_URL", Required: true},
	}
	violations := Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("expected key DATABASE_URL, got %s", violations[0].Key)
	}
}

func TestValidate_EmptyValueRequired(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "   "}
	rules := []Rule{
		{Key: "SECRET_KEY", Required: true},
	}
	violations := Validate(env, rules)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation for blank value, got %d", len(violations))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{"APP_PORT": "not-a-port"}
	rules := []Rule{
		{Key: "APP_PORT", Required: false, Pattern: `^\d+$`},
	}
	violations := Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_InvalidPattern(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	rules := []Rule{
		{Key: "FOO", Pattern: `[invalid(`},
	}
	violations := Validate(env, rules)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation for invalid regex, got %d", len(violations))
	}
}

func TestValidate_OptionalMissingNoViolation(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{
		{Key: "OPTIONAL_VAR", Required: false, Pattern: `^\d+$`},
	}
	violations := Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for optional missing key, got %v", violations)
	}
}
