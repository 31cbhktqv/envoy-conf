package resolver

import (
	"fmt"
	"os"
	"strings"
)

// Source represents a named environment variable source.
type Source struct {
	Name string
	Vars map[string]string
}

// ResolveOptions controls how environment variable resolution behaves.
type ResolveOptions struct {
	// FallbackToOS will merge OS environment variables as a base layer.
	FallbackToOS bool
	// OverrideKeys lists keys whose values should be overridden by OS env.
	OverrideKeys []string
}

// Resolve merges multiple Sources into a single flat map.
// Sources are applied in order; later sources override earlier ones.
// If opts.FallbackToOS is true, OS environment is used as the base layer.
func Resolve(sources []Source, opts ResolveOptions) (map[string]string, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("resolver: at least one source is required")
	}

	result := make(map[string]string)

	if opts.FallbackToOS {
		for _, entry := range os.Environ() {
			parts := strings.SplitN(entry, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	}

	for _, src := range sources {
		for k, v := range src.Vars {
			result[k] = v
		}
	}

	for _, key := range opts.OverrideKeys {
		if val, ok := os.LookupEnv(key); ok {
			result[key] = val
		}
	}

	return result, nil
}

// SourceNames returns the names of all provided sources.
func SourceNames(sources []Source) []string {
	names := make([]string, len(sources))
	for i, s := range sources {
		names[i] = s.Name
	}
	return names
}
