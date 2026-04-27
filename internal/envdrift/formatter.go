package envdrift

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// RenderOptions controls output behaviour.
type RenderOptions struct {
	Color      bool
	ShowMatch  bool
}

// Render writes a human-readable drift report to w.
func Render(w io.Writer, r Report, opts RenderOptions) {
	fmt.Fprintf(w, "Drift report for target: %s\n", r.Target)
	fmt.Fprintln(w, strings.Repeat("-", 48))

	for _, e := range r.Entries {
		switch e.Status {
		case StatusMatch:
			if opts.ShowMatch {
				fmt.Fprintf(w, "  %s\n", e.Key)
			}
		case StatusAdded:
			line := fmt.Sprintf("+ %s=%s", e.Key, e.Live)
			if opts.Color {
				line = colorGreen + line + colorReset
			}
			fmt.Fprintln(w, line)
		case StatusRemoved:
			line := fmt.Sprintf("- %s=%s", e.Key, e.Baseline)
			if opts.Color {
				line = colorRed + line + colorReset
			}
			fmt.Fprintln(w, line)
		case StatusChanged:
			line := fmt.Sprintf("~ %s: %q → %q", e.Key, e.Baseline, e.Live)
			if opts.Color {
				line = colorYellow + line + colorReset
			}
			fmt.Fprintln(w, line)
		}
	}

	added, removed, changed := r.Counts()
	fmt.Fprintln(w, strings.Repeat("-", 48))
	fmt.Fprintf(w, "Summary: +%d added  -%d removed  ~%d changed\n", added, removed, changed)
}
