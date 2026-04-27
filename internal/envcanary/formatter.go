package envcanary

import (
	"fmt"
	"io"
	"strings"
)

const (
	iconOK       = "✔"
	iconWarning  = "⚠"
	iconCritical = "✘"
)

// Render writes a human-readable canary report to w.
func Render(w io.Writer, results []Result, showOK bool) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No canary checks configured.")
		return
	}

	for _, r := range results {
		if r.Status == StatusOK && !showOK {
			continue
		}
		icon := statusIcon(r.Status)
		fmt.Fprintf(w, "  %s  %-30s %s\n", icon, r.Key, r.Message)
		if r.Status != StatusOK && r.Baseline != "" {
			fmt.Fprintf(w, "       baseline: %s\n", r.Baseline)
			fmt.Fprintf(w, "       current:  %s\n", r.Current)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintln(w, RenderSummary(results))
}

// RenderSummary returns a one-line summary string.
func RenderSummary(results []Result) string {
	ok, warn, crit := 0, 0, 0
	for _, r := range results {
		switch r.Status {
		case StatusOK:
			ok++
		case StatusWarning:
			warn++
		case StatusCritical:
			crit++
		}
	}
	return fmt.Sprintf("canary: %d ok, %d warning, %d critical", ok, warn, crit)
}

func statusIcon(s Status) string {
	switch s {
	case StatusOK:
		return iconOK
	case StatusWarning:
		return iconWarning
	case StatusCritical:
		return iconCritical
	default:
		return "?"
	}
}
