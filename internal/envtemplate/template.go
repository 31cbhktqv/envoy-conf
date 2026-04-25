package envtemplate

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} or $VAR_NAME style placeholders.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls template expansion behaviour.
type Options struct {
	// Strict causes Expand to return an error if any variable is unresolved.
	Strict bool
	// Fallback is returned for unresolved variables when Strict is false.
	Fallback string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Strict: false, Fallback: ""}
}

// Expand replaces variable placeholders in each value of env using the
// provided lookup map. OS environment variables are also consulted as a
// secondary source.
func Expand(env map[string]string, lookup map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := expandValue(v, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

func expandValue(value string, lookup map[string]string, opts Options) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := lookup[name]; ok {
			return v
		}
		if v, ok := os.LookupEnv(name); ok {
			return v
		}
		if opts.Strict {
			expandErr = fmt.Errorf("unresolved variable: %s", name)
			return match
		}
		return opts.Fallback
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
