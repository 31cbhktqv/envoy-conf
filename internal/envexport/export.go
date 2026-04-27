package envexport

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents the output format for exported env vars.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport  Format = "export"
	FormatJSON    Format = "json"
	FormatYAML    Format = "yaml"
)

// Options controls how env vars are exported.
type Options struct {
	Format  Format
	Sorted  bool
	Masked  map[string]bool
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatDotenv,
		Sorted: true,
	}
}

// Export serializes the given env map to the requested format.
func Export(env map[string]string, opts Options) (string, error) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	switch opts.Format {
	case FormatDotenv:
		return renderDotenv(keys, env, opts), nil
	case FormatExport:
		return renderExport(keys, env, opts), nil
	case FormatJSON:
		return renderJSON(keys, env, opts), nil
	case FormatYAML:
		return renderYAML(keys, env, opts), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func masked(key string, opts Options) bool {
	return opts.Masked != nil && opts.Masked[key]
}

func val(key, v string, opts Options) string {
	if masked(key, opts) {
		return "***"
	}
	return v
}

func renderDotenv(keys []string, env map[string]string, opts Options) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, val(k, env[k], opts))
	}
	return sb.String()
}

func renderExport(keys []string, env map[string]string, opts Options) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, val(k, env[k], opts))
	}
	return sb.String()
}

func renderJSON(keys []string, env map[string]string, opts Options) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, val(k, env[k], opts), comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func renderYAML(keys []string, env map[string]string, opts Options) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s: %q\n", k, val(k, env[k], opts))
	}
	return sb.String()
}
