package envschema

import (
	"testing"
)

func schema() *Schema {
	return &Schema{
		Fields: []Field{
			{Key: "APP_ENV", Type: TypeString, Required: true, Allowed: []string{"development", "staging", "production"}},
			{Key: "PORT", Type: TypeInt, Required: true},
			{Key: "DEBUG", Type: TypeBool, Required: false},
			{Key: "API_URL", Type: TypeURL, Required: false},
			{Key: "LOG_LEVEL", Type: TypeString, Required: false, Pattern: `^(debug|info|warn|error)$`},
		},
	}
}

func TestValidate_AllPass(t *testing.T) {
	env := map[string]string{
		"APP_ENV":   "production",
		"PORT":      "8080",
		"DEBUG":     "false",
		"API_URL":   "https://api.example.com",
		"LOG_LEVEL": "info",
	}
	violations := schema().Validate(env)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging"}
	violations := schema().Validate(env)
	if len(violations) != 1 || violations[0].Key != "PORT" {
		t.Fatalf("expected PORT violation, got %v", violations)
	}
}

func TestValidate_InvalidType_Int(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging", "PORT": "not-a-number"}
	violations := schema().Validate(env)
	if len(violations) != 1 || violations[0].Key != "PORT" {
		t.Fatalf("expected PORT type violation, got %v", violations)
	}
}

func TestValidate_InvalidType_Bool(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging", "PORT": "3000", "DEBUG": "maybe"}
	violations := schema().Validate(env)
	if len(violations) != 1 || violations[0].Key != "DEBUG" {
		t.Fatalf("expected DEBUG violation, got %v", violations)
	}
}

func TestValidate_InvalidType_URL(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging", "PORT": "3000", "API_URL": "not-a-url"}
	violations := schema().Validate(env)
	if len(violations) != 1 || violations[0].Key != "API_URL" {
		t.Fatalf("expected API_URL violation, got %v", violations)
	}
}

func TestValidate_AllowedValues(t *testing.T) {
	env := map[string]string{"APP_ENV": "local", "PORT": "3000"}
	violations := schema().Validate(env)
	if len(violations) != 1 || violations[0].Key != "APP_ENV" {
		t.Fatalf("expected APP_ENV allowlist violation, got %v", violations)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging", "PORT": "3000", "LOG_LEVEL": "verbose"}
	violations := schema().Validate(env)
	if len(violations) != 1 || violations[0].Key != "LOG_LEVEL" {
		t.Fatalf("expected LOG_LEVEL pattern violation, got %v", violations)
	}
}
