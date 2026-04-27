package envpromote

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Render writes a human-readable promotion report to w.
func Render(w io.Writer, r PromoteResult) {
	fmt.Fprintf(w, "Promotion: %s → %s\n", r.From, r.To)
	fmt.Fprintf(w, "Status:    %s\n\n", readyLabel(r.Ready))

	if len(r.MissingKeys) > 0 {
		sort.Strings(r.MissingKeys)
		fmt.Fprintln(w, "Missing keys (blocking):")
		for _, k := range r.MissingKeys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
		fmt.Fprintln(w)
	}

	if len(r.ChangedKeys) > 0 {
		sort.Strings(r.ChangedKeys)
		fmt.Fprintln(w, "Changed values (informational):")
		for _, k := range r.ChangedKeys {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
		fmt.Fprintln(w)
	}

	if len(r.NewKeys) > 0 {
		sort.Strings(r.NewKeys)
		fmt.Fprintln(w, "New keys in target (informational):")
		for _, k := range r.NewKeys {
			fmt.Fprintf(w, "  + %s\n", k)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintln(w, Summary(r))
}

func readyLabel(ready bool) string {
	if ready {
		return "READY"
	}
	return "NOT READY"
}
