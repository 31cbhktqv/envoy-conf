package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/envoy-conf/internal/differ"
	"github.com/user/envoy-conf/internal/validator"
)

func baseReport() Report {
	return Report{
		GeneratedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		SourceFile:  "prod.env",
		TargetFile:  "staging.env",
	}
}

func TestWrite_TextFormat_NoDifferences(t *testing.T) {
	r := baseReport()
	r.DiffSummary = "No differences found."

	var buf bytes.Buffer
	if err := Write(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected source file in output")
	}
	if !strings.Contains(out, "all checks passed") {
		t.Errorf("expected validation pass message")
	}
}

func TestWrite_TextFormat_WithViolations(t *testing.T) {
	r := baseReport()
	r.Violations = []validator.Violation{
		{Key: "DB_HOST", Message: "required key is missing", Severity: "error"},
	}

	var buf bytes.Buffer
	if err := Write(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "1 violation(s)") {
		t.Errorf("expected violation count in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected violation key in output")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	r := baseReport()
	r.DiffItems = []differ.DiffItem{
		{Key: "PORT", Status: differ.StatusChanged, ValueA: "8080", ValueB: "9090"},
	}

	var buf bytes.Buffer
	if err := Write(&buf, r, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded.DiffItems) != 1 {
		t.Errorf("expected 1 diff item, got %d", len(decoded.DiffItems))
	}
	if decoded.SourceFile != "prod.env" {
		t.Errorf("expected source file in JSON output")
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, baseReport(), Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
