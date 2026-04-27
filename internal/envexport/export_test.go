package envexport

import (
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_PASS":  "secret123",
	"LOG_LEVEL": "info",
}

func TestExport_Dotenv(t *testing.T) {
	opts := DefaultOptions()
	out, err := Export(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line, got:\n%s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL=info") {
		t.Errorf("expected LOG_LEVEL line, got:\n%s", out)
	}
}

func TestExport_ExportFormat(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatExport
	out, err := Export(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=\"") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatJSON
	out, err := Export(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV": "production"`) {
		t.Errorf("expected JSON key/value, got:\n%s", out)
	}
}

func TestExport_YAMLFormat(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatYAML
	out, err := Export(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV: ") {
		t.Errorf("expected YAML line, got:\n%s", out)
	}
}

func TestExport_MaskedKeys(t *testing.T) {
	opts := DefaultOptions()
	opts.Masked = map[string]bool{"DB_PASS": true}
	out, err := Export(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "secret123") {
		t.Errorf("expected masked value, got raw secret in:\n%s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected *** mask in output:\n%s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = Format("xml")
	_, err := Export(sampleEnv, opts)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExport_SortedOutput(t *testing.T) {
	opts := DefaultOptions()
	out, err := Export(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected APP_ENV first (sorted), got: %s", lines[0])
	}
}
