package linter

import (
	"fmt"
	"strings"
)

// Rule defines a linting rule applied to environment variable keys or values.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Violation represents a single linting issue found in an env map.
type Violation struct {
	Key     string
	Rule    string
	Message string
}

// Linter holds a set of rules to apply against an env map.
type Linter struct {
	rules []Rule
}

// DefaultRules returns a standard set of linting rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "no-lowercase-key",
			Message: "environment variable key should be uppercase",
			Check: func(key, _ string) bool {
				return key != strings.ToUpper(key)
			},
		},
		{
			Name:    "no-empty-value",
			Message: "environment variable has an empty value",
			Check: func(_, value string) bool {
				return strings.TrimSpace(value) == ""
			},
		},
		{
			Name:    "no-whitespace-in-key",
			Message: "environment variable key contains whitespace",
			Check: func(key, _ string) bool {
				return strings.ContainsAny(key, " \t")
			},
		},
	}
}

// New creates a Linter with the provided rules.
func New(rules []Rule) *Linter {
	return &Linter{rules: rules}
}

// Lint runs all rules against the provided env map and returns any violations.
func (l *Linter) Lint(env map[string]string) []Violation {
	var violations []Violation
	for key, value := range env {
		for _, rule := range l.rules {
			if rule.Check(key, value) {
				violations = append(violations, Violation{
					Key:     key,
					Rule:    rule.Name,
					Message: fmt.Sprintf("%s: %s", key, rule.Message),
				})
			}
		}
	}
	return violations
}
