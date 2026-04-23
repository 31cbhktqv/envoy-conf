package differ

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

// FormatOptions controls how the diff output is rendered.
type FormatOptions struct {
	Color   bool
	Verbose bool
	LabelA  string
	LabelB  string
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Color:  true,
		LabelA: "source",
		LabelB: "target",
	}
}

// Render writes a human-readable diff to w.
func Render(w io.Writer, d DiffResult, opts FormatOptions) {
	if !d.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	for _, k := range d.SortedOnlyInAKeys() {
		line := fmt.Sprintf("- %s=%s  (only in %s)", k, d.OnlyInA[k], opts.LabelA)
		fmt.Fprintln(w, colorize(line, colorRed, opts.Color))
	}

	for _, k := range d.SortedOnlyInBKeys() {
		line := fmt.Sprintf("+ %s=%s  (only in %s)", k, d.OnlyInB[k], opts.LabelB)
		fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))
	}

	for _, k := range d.SortedChangedKeys() {
		pair := d.Changed[k]
		line := fmt.Sprintf("~ %s: %s → %s", k, pair[0], pair[1])
		fmt.Fprintln(w, colorize(line, colorYellow, opts.Color))
	}

	if opts.Verbose {
		for k, v := range d.Unchanged {
			fmt.Fprintf(w, "  %s=%s\n", k, v)
		}
	}
}

// Summary returns a one-line summary string.
func Summary(d DiffResult) string {
	parts := []string{}
	if n := len(d.OnlyInA); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(d.OnlyInB); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(d.Changed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}
