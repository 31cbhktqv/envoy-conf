package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/envoy-conf/internal/differ"
	"github.com/user/envoy-conf/internal/validator"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the combined diff and validation results.
type Report struct {
	GeneratedAt time.Time              `json:"generated_at"`
	SourceFile  string                 `json:"source_file"`
	TargetFile  string                 `json:"target_file"`
	DiffItems   []differ.DiffItem      `json:"diff,omitempty"`
	Violations  []validator.Violation  `json:"violations,omitempty"`
	DiffSummary string                 `json:"diff_summary,omitempty"`
}

// Write renders the report in the specified format to the given writer.
func Write(w io.Writer, r Report, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, r)
	case FormatText:
		return writeText(w, r)
	default:
		return fmt.Errorf("unsupported report format: %q", format)
	}
}

func writeJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func writeText(w io.Writer, r Report) error {
	fmt.Fprintf(w, "Report generated: %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Source : %s\n", r.SourceFile)
	fmt.Fprintf(w, "Target : %s\n\n", r.TargetFile)

	if r.DiffSummary != "" {
		fmt.Fprintf(w, "Diff Summary:\n%s\n", r.DiffSummary)
	}

	if len(r.Violations) == 0 {
		fmt.Fprintln(w, "Validation: all checks passed")
	} else {
		fmt.Fprintf(w, "Validation: %d violation(s) found\n", len(r.Violations))
		for _, v := range r.Violations {
			fmt.Fprintf(w, "  [%s] %s: %s\n", v.Severity, v.Key, v.Message)
		}
	}
	return nil
}

// WriteFile writes the report to a file at the given path.
func WriteFile(path string, r Report, format Format) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating report file: %w", err)
	}
	defer f.Close()
	return Write(f, r, format)
}
