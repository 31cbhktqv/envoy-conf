package envcanary

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender_NoChecks(t *testing.T) {
	var buf bytes.Buffer
	Render(&buf, nil, true)
	if !strings.Contains(buf.String(), "No canary checks configured") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestRender_ShowOK(t *testing.T) {
	results := []Result{
		{Key: "APP_ENV", Status: StatusOK, Current: "prod", Message: "present"},
	}
	var buf bytes.Buffer
	Render(&buf, results, true)
	if !strings.Contains(buf.String(), iconOK) {
		t.Errorf("expected OK icon in output: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "APP_ENV") {
		t.Errorf("expected key name in output")
	}
}

func TestRender_HideOK(t *testing.T) {
	results := []Result{
		{Key: "APP_ENV", Status: StatusOK, Current: "prod", Message: "present"},
	}
	var buf bytes.Buffer
	Render(&buf, results, false)
	if strings.Contains(buf.String(), "APP_ENV") {
		t.Errorf("OK result should be hidden when showOK=false")
	}
}

func TestRender_Critical(t *testing.T) {
	results := []Result{
		{Key: "DB_URL", Status: StatusCritical, Baseline: "postgres://old", Current: "", Message: "required key missing"},
	}
	var buf bytes.Buffer
	Render(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, iconCritical) {
		t.Errorf("expected critical icon")
	}
	if !strings.Contains(out, "postgres://old") {
		t.Errorf("expected baseline value in output")
	}
}

func TestRenderSummary(t *testing.T) {
	results := []Result{
		{Status: StatusOK},
		{Status: StatusWarning},
		{Status: StatusCritical},
		{Status: StatusCritical},
	}
	summary := RenderSummary(results)
	if !strings.Contains(summary, "1 ok") {
		t.Errorf("expected 1 ok in summary: %q", summary)
	}
	if !strings.Contains(summary, "1 warning") {
		t.Errorf("expected 1 warning in summary: %q", summary)
	}
	if !strings.Contains(summary, "2 critical") {
		t.Errorf("expected 2 critical in summary: %q", summary)
	}
}
