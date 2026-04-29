package envgraph

import (
	"testing"
)

func TestBuild_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	g := Build(env)
	if len(g.nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(g.nodes))
	}
	for _, n := range g.nodes {
		if len(n.Deps) != 0 {
			t.Errorf("expected no deps for %q", n.Key)
		}
	}
}

func TestBuild_WithRefs(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "localhost",
		"PORT":     "8080",
	}
	g := Build(env)
	n := g.nodes["BASE_URL"]
	if n == nil {
		t.Fatal("BASE_URL node missing")
	}
	if len(n.Deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(n.Deps), n.Deps)
	}
}

func TestOrder_NoCycle(t *testing.T) {
	g := New()
	g.Add("A", []string{})
	g.Add("B", []string{"A"})
	g.Add("C", []string{"B"})

	order, err := g.Order()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := func(k string) int {
		for i, v := range order {
			if v == k {
				return i
			}
		}
		return -1
	}
	if pos("A") > pos("B") || pos("B") > pos("C") {
		t.Errorf("wrong order: %v", order)
	}
}

func TestOrder_CycleDetected(t *testing.T) {
	g := New()
	g.Add("A", []string{"B"})
	g.Add("B", []string{"A"})

	_, err := g.Order()
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestExtractRefs_CurlyBrace(t *testing.T) {
	refs := extractRefs("http://${HOST}:${PORT}/path")
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %v", refs)
	}
}

func TestExtractRefs_BareVar(t *testing.T) {
	refs := extractRefs("$SCHEME://$HOST")
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %v", refs)
	}
}

func TestExtractRefs_NoDuplicates(t *testing.T) {
	refs := extractRefs("${X}-${X}-${X}")
	if len(refs) != 1 {
		t.Fatalf("expected 1 unique ref, got %v", refs)
	}
}
