package envgraph

import (
	"fmt"
	"io"
	"strings"
)

// RenderOptions controls output of Render.
type RenderOptions struct {
	ShowDeps bool
}

// DefaultRenderOptions returns sensible defaults.
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{ShowDeps: true}
}

// Render writes a human-readable dependency graph to w.
func Render(w io.Writer, g *Graph, opts RenderOptions) error {
	order, err := g.Order()
	if err != nil {
		return fmt.Errorf("cannot render graph: %w", err)
	}

	if len(order) == 0 {
		fmt.Fprintln(w, "(no variables)")
		return nil
	}

	for _, k := range order {
		n, ok := g.nodes[k]
		if !ok {
			continue
		}
		if opts.ShowDeps && len(n.Deps) > 0 {
			fmt.Fprintf(w, "%s -> [%s]\n", k, strings.Join(n.Deps, ", "))
		} else {
			fmt.Fprintln(w, k)
		}
	}
	return nil
}

// RenderDOT writes a Graphviz DOT representation to w.
func RenderDOT(w io.Writer, g *Graph) error {
	order, err := g.Order()
	if err != nil {
		return fmt.Errorf("cannot render DOT: %w", err)
	}

	fmt.Fprintln(w, "digraph envgraph {")
	fmt.Fprintln(w, `  rankdir=LR;`)

	for _, k := range order {
		n, ok := g.nodes[k]
		if !ok {
			continue
		}
		if len(n.Deps) == 0 {
			fmt.Fprintf(w, "  %q;\n", k)
		}
		for _, dep := range n.Deps {
			fmt.Fprintf(w, "  %q -> %q;\n", dep, k)
		}
	}

	fmt.Fprintln(w, "}")
	return nil
}
