package masker

import (
	"regexp"
	"strings"
)

// DefaultSensitivePatterns are common patterns for sensitive environment variable names.
var DefaultSensitivePatterns = []string{
	"(?i)password",
	"(?i)secret",
	"(?i)token",
	"(?i)api_key",
	"(?i)apikey",
	"(?i)private_key",
	"(?i)auth",
	"(?i)credential",
}

// Masker redacts sensitive environment variable values.
type Masker struct {
	patterns []*regexp.Regexp
	mask     string
}

// New creates a Masker with the given key patterns and mask string.
// If patterns is nil, DefaultSensitivePatterns are used.
// If mask is empty, "***" is used.
func New(patterns []string, mask string) (*Masker, error) {
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}
	if mask == "" {
		mask = "***"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Masker{patterns: compiled, mask: mask}, nil
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, re := range m.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// MaskEnv returns a copy of the env map with sensitive values redacted.
func (m *Masker) MaskEnv(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if m.IsSensitive(k) {
			out[k] = m.mask
		} else {
			out[k] = v
		}
	}
	return out
}

// MaskValue redacts a single value if the key is sensitive, otherwise returns it unchanged.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.mask
	}
	return strings.Clone(value)
}
