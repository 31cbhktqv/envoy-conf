package envcompare

import (
	"fmt"
	"regexp"
)

// CompareOptions controls how two env maps are compared.
type CompareOptions struct {
	// IgnoreKeys is a list of exact key names to skip during comparison.
	IgnoreKeys []string
	// IgnorePatterns is a list of regex patterns; matching keys are skipped.
	IgnorePatterns []string
	// CaseSensitiveValues controls whether value comparison is case-sensitive.
	CaseSensitiveValues bool
}

// Result holds the outcome of a comparison between two env maps.
type Result struct {
	MissingInB  []string          // keys present in A but not in B
	MissingInA  []string          // keys present in B but not in A
	Mismatched  map[string][2]string // key -> [valueA, valueB]
	MatchedCount int
}

// Compare compares two env maps according to the provided options.
func Compare(a, b map[string]string, opts CompareOptions) (Result, error) {
	ignoreSet := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignoreSet[k] = true
	}

	patterns := make([]*regexp.Regexp, 0, len(opts.IgnorePatterns))
	for _, p := range opts.IgnorePatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return Result{}, fmt.Errorf("invalid ignore pattern %q: %w", p, err)
		}
		patterns = append(patterns, re)
	}

	shouldIgnore := func(key string) bool {
		if ignoreSet[key] {
			return true
		}
		for _, re := range patterns {
			if re.MatchString(key) {
				return true
			}
		}
		return false
	}

	res := Result{
		Mismatched: make(map[string][2]string),
	}

	for k, va := range a {
		if shouldIgnore(k) {
			continue
		}
		vb, ok := b[k]
		if !ok {
			res.MissingInB = append(res.MissingInB, k)
			continue
		}
		equal := va == vb
		if !opts.CaseSensitiveValues {
			equal = stringsEqualFold(va, vb)
		}
		if !equal {
			res.Mismatched[k] = [2]string{va, vb}
		} else {
			res.MatchedCount++
		}
	}

	for k := range b {
		if shouldIgnore(k) {
			continue
		}
		if _, ok := a[k]; !ok {
			res.MissingInA = append(res.MissingInA, k)
		}
	}

	return res, nil
}

func stringsEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
