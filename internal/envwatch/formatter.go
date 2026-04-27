package envwatch

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

// RenderOptions controls output formatting.
type RenderOptions struct {
	Color   bool
	Verbose bool
}

// Render writes a human-readable summary of changes to w.
func Render(w io.Writer, changes []Change, opts RenderOptions) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "no changes detected")
		return
	}

	sorted := make([]Change, len(changes))
	copy(sorted, changes)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	for _, c := range sorted {
		switch c.Type {
		case Added:
			line := fmt.Sprintf("+ %s=%s", c.Key, c.NewVal)
			fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))
		case Removed:
			line := fmt.Sprintf("- %s=%s", c.Key, c.OldVal)
			fmt.Fprintln(w, colorize(line, colorRed, opts.Color))
		case Changed:
			if opts.Verbose {
				line := fmt.Sprintf("~ %s: %s → %s", c.Key, c.OldVal, c.NewVal)
				fmt.Fprintln(w, colorize(line, colorYellow, opts.Color))
			} else {
				line := fmt.Sprintf("~ %s=%s", c.Key, c.NewVal)
				fmt.Fprintln(w, colorize(line, colorYellow, opts.Color))
			}
		}
	}
}

// RenderSummary returns a one-line summary string.
func RenderSummary(changes []Change) string {
	var added, removed, changed int
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}
