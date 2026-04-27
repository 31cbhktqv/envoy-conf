package envwatch

import (
	"strings"
	"testing"
)

func TestRender_NoChanges(t *testing.T) {
	var sb strings.Builder
	Render(&sb, nil, RenderOptions{})
	if !strings.Contains(sb.String(), "no changes") {
		t.Fatalf("expected no changes message, got: %q", sb.String())
	}
}

func TestRender_Added(t *testing.T) {
	var sb strings.Builder
	Render(&sb, []Change{{Key: "FOO", Type: Added, NewVal: "bar"}}, RenderOptions{})
	out := sb.String()
	if !strings.Contains(out, "+ FOO=bar") {
		t.Fatalf("expected added line, got: %q", out)
	}
}

func TestRender_Removed(t *testing.T) {
	var sb strings.Builder
	Render(&sb, []Change{{Key: "FOO", Type: Removed, OldVal: "bar"}}, RenderOptions{})
	out := sb.String()
	if !strings.Contains(out, "- FOO=bar") {
		t.Fatalf("expected removed line, got: %q", out)
	}
}

func TestRender_Changed_Verbose(t *testing.T) {
	var sb strings.Builder
	Render(&sb, []Change{{Key: "K", Type: Changed, OldVal: "old", NewVal: "new"}}, RenderOptions{Verbose: true})
	out := sb.String()
	if !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Fatalf("expected verbose change line, got: %q", out)
	}
}

func TestRender_SortedOutput(t *testing.T) {
	var sb strings.Builder
	Render(&sb, []Change{
		{Key: "Z", Type: Added, NewVal: "1"},
		{Key: "A", Type: Added, NewVal: "2"},
	}, RenderOptions{})
	out := sb.String()
	if strings.Index(out, "A") > strings.Index(out, "Z") {
		t.Fatal("expected sorted output")
	}
}

func TestRenderSummary(t *testing.T) {
	changes := []Change{
		{Key: "A", Type: Added},
		{Key: "B", Type: Removed},
		{Key: "C", Type: Changed},
	}
	s := RenderSummary(changes)
	for _, want := range []string{"1 added", "1 removed", "1 changed"} {
		if !strings.Contains(s, want) {
			t.Fatalf("expected %q in summary %q", want, s)
		}
	}
}

func TestRenderSummary_Empty(t *testing.T) {
	if RenderSummary(nil) != "no changes" {
		t.Fatal("expected 'no changes'")
	}
}
