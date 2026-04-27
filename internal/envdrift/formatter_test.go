package envdrift

import (
	"strings"
	"testing"
)

func TestRender_NoDrift(t *testing.T) {
	r := Detect("prod", map[string]string{"A": "1"}, map[string]string{"A": "1"})
	var sb strings.Builder
	Render(&sb, r, RenderOptions{})
	out := sb.String()
	if !strings.Contains(out, "prod") {
		t.Error("expected target name in output")
	}
	if !strings.Contains(out, "+0 added") {
		t.Errorf("expected zero summary, got: %s", out)
	}
}

func TestRender_ShowsAdded(t *testing.T) {
	r := Detect("staging", map[string]string{}, map[string]string{"NEW": "val"})
	var sb strings.Builder
	Render(&sb, r, RenderOptions{})
	if !strings.Contains(sb.String(), "+ NEW=val") {
		t.Errorf("expected added line, got: %s", sb.String())
	}
}

func TestRender_ShowsRemoved(t *testing.T) {
	r := Detect("staging", map[string]string{"OLD": "gone"}, map[string]string{})
	var sb strings.Builder
	Render(&sb, r, RenderOptions{})
	if !strings.Contains(sb.String(), "- OLD=gone") {
		t.Errorf("expected removed line, got: %s", sb.String())
	}
}

func TestRender_ShowsChanged(t *testing.T) {
	r := Detect("prod", map[string]string{"HOST": "old"}, map[string]string{"HOST": "new"})
	var sb strings.Builder
	Render(&sb, r, RenderOptions{})
	if !strings.Contains(sb.String(), "~ HOST") {
		t.Errorf("expected changed line, got: %s", sb.String())
	}
}

func TestRender_ShowMatch_Option(t *testing.T) {
	r := Detect("prod", map[string]string{"SAME": "x"}, map[string]string{"SAME": "x"})
	var sb strings.Builder
	Render(&sb, r, RenderOptions{ShowMatch: true})
	if !strings.Contains(sb.String(), "SAME") {
		t.Errorf("expected match key in output when ShowMatch=true")
	}
}
