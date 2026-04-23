package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks the provided env map against the given rules.
// It returns a slice of Violations (empty means all rules passed).
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && (!exists || strings.TrimSpace(val) == "") {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required but missing or empty",
			})
			continue
		}

		if exists && rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}
	}

	return violations
}
