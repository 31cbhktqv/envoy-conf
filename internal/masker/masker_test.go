package masker

import (
	"testing"
)

func TestNew_DefaultPatterns(t *testing.T) {
	m, err := New(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
	if m.mask != "***" {
		t.Errorf("expected default mask '***', got %q", m.mask)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]string{"[invalid"}, "")
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestIsSensitive(t *testing.T) {
	m, _ := New(nil, "")

	sensitiveKeys := []string{
		"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN",
		"STRIPE_API_KEY", "PRIVATE_KEY", "APP_CREDENTIAL",
	}
	for _, k := range sensitiveKeys {
		if !m.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}

	safeKeys := []string{"PORT", "HOST", "LOG_LEVEL", "APP_ENV"}
	for _, k := range safeKeys {
		if m.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestMaskEnv(t *testing.T) {
	m, _ := New(nil, "REDACTED")

	env := map[string]string{
		"PORT":        "8080",
		"DB_PASSWORD": "supersecret",
		"APP_ENV":     "production",
		"API_TOKEN":   "tok_abc123",
	}

	masked := m.MaskEnv(env)

	if masked["PORT"] != "8080" {
		t.Errorf("expected PORT to be unmasked, got %q", masked["PORT"])
	}
	if masked["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unmasked, got %q", masked["APP_ENV"])
	}
	if masked["DB_PASSWORD"] != "REDACTED" {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["API_TOKEN"] != "REDACTED" {
		t.Errorf("expected API_TOKEN to be masked, got %q", masked["API_TOKEN"])
	}

	// ensure original is not modified
	if env["DB_PASSWORD"] != "supersecret" {
		t.Error("original env map should not be modified")
	}
}

func TestMaskValue(t *testing.T) {
	m, _ := New(nil, "")

	if got := m.MaskValue("DB_PASSWORD", "secret"); got != "***" {
		t.Errorf("expected masked value, got %q", got)
	}
	if got := m.MaskValue("PORT", "8080"); got != "8080" {
		t.Errorf("expected plain value, got %q", got)
	}
}
