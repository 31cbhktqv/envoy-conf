package envschema

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldType represents the expected type of an environment variable value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeURL     FieldType = "url"
)

// Field describes the schema for a single environment variable.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  string // optional regex
	Allowed  []string // optional allowlist
}

// Schema holds a collection of field definitions.
type Schema struct {
	Fields []Field
}

// Violation describes a schema validation failure.
type Violation struct {
	Key     string
	Message string
}

var (
	reInt  = regexp.MustCompile(`^-?\d+$`)
	reBool = regexp.MustCompile(`^(?i)(true|false|1|0|yes|no)$`)
	reURL  = regexp.MustCompile(`^https?://`)
)

// Validate checks env against the schema and returns any violations.
func (s *Schema) Validate(env map[string]string) []Violation {
	var violations []Violation

	for _, f := range s.Fields {
		val, exists := env[f.Key]

		if !exists || strings.TrimSpace(val) == "" {
			if f.Required {
				violations = append(violations, Violation{Key: f.Key, Message: "required key is missing or empty"})
			}
			continue
		}

		switch f.Type {
		case TypeInt:
			if !reInt.MatchString(val) {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("expected int, got %q", val)})
			}
		case TypeBool:
			if !reBool.MatchString(val) {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("expected bool, got %q", val)})
			}
		case TypeURL:
			if !reURL.MatchString(val) {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("expected URL, got %q", val)})
			}
		}

		if f.Pattern != "" {
			re, err := regexp.Compile(f.Pattern)
			if err == nil && !re.MatchString(val) {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("value %q does not match pattern %q", val, f.Pattern)})
			}
		}

		if len(f.Allowed) > 0 {
			matched := false
			for _, a := range f.Allowed {
				if a == val {
					matched = true
					break
				}
			}
			if !matched {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("value %q not in allowed list %v", val, f.Allowed)})
			}
		}
	}

	return violations
}
