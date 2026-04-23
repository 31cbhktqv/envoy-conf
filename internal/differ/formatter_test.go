package differ

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender_NoDifferences(t *testing.T) {
	d := Diff(
		map[string]string{"FOO": "bar"},
		map[string]string{"FOO": "bar"},
	)
	var buf bytes.Buffer
	Render(&buf, d, FormatOptions{Color: false, LabelA: "a", LabelB: "b"})
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestRender_ShowsAdded(t *testing.T) {
	d := Diff(
		map[string]string{},
		map[string]string{"NEW_KEY": "value"},
	)
	var buf bytes.Buffer
	Render(&buf, d, FormatOptions{Color: false, LabelA: "src", LabelB: "tgt"})
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY=value") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestRender_ShowsRemoved(t *testing.T) {
	d := Diff(
		map[string]string{"OLD_KEY": "val"},
		map[string]string{},
	)
	var buf bytes.Buffer
	Render(&buf, d, FormatOptions{Color: false, LabelA: "src", LabelB: "tgt"})
	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY=val") {
		t.Errorf("expected removed key in output, got: %s", out)
	}
}

func TestRender_ShowsChanged(t *testing.T) {
	d := Diff(
		map[string]string{"PORT": "8080"},
		map[string]string{"PORT": "9090"},
	)
	var buf bytes.Buffer
	Render(&buf, d, FormatOptions{Color: false})
	out := buf.String()
	if !strings.Contains(out, "~ PORT") {
		t.Errorf("expected changed key in output, got: %s", out)
	}
}

func TestSummary(t *testing.T) {
	d := Diff(
		map[string]string{"A": "1", "B": "old"},
		map[string]string{"B": "new", "C": "3"},
	)
	s := Summary(d)
	if !strings.Contains(s, "removed") || !strings.Contains(s, "added") || !strings.Contains(s, "changed") {
		t.Errorf("unexpected summary: %s", s)
	}
}
