package envfilter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// Prefix keeps only keys with the given prefix (case-insensitive).
	Prefix string
	// Pattern keeps only keys matching the regular expression.
	Pattern string
	// ExcludeKeys is a set of exact key names to drop.
	ExcludeKeys []string
}

// Filter returns a new map containing only the entries from env that
// satisfy all criteria defined in opts.
func Filter(env map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	exclude := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range env {
		if _, skip := exclude[k]; skip {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(strings.ToUpper(k), strings.ToUpper(opts.Prefix)) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		result[k] = v
	}
	return result, nil
}
