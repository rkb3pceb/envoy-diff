package env

import (
	"fmt"
	"strconv"
	"strings"
)

// SchemaType represents the expected type of an environment variable value.
type SchemaType string

const (
	SchemaTypeString SchemaType = "string"
	SchemaTypeInt    SchemaType = "int"
	SchemaTypeBool   SchemaType = "bool"
	SchemaTypeURL    SchemaType = "url"
)

// SchemaRule defines a validation rule for a single env key.
type SchemaRule struct {
	Key      string
	Type     SchemaType
	Required bool
}

// SchemaFinding is a result of schema validation.
type SchemaFinding struct {
	Key     string
	Message string
	Error   bool
}

// ValidateSchema checks env map values against a list of schema rules.
// It returns a slice of findings (type mismatches, missing required keys).
func ValidateSchema(env map[string]string, rules []SchemaRule) []SchemaFinding {
	var findings []SchemaFinding

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if !exists || val == "" {
			if rule.Required {
				findings = append(findings, SchemaFinding{
					Key:     rule.Key,
					Message: fmt.Sprintf("required key %q is missing or empty", rule.Key),
					Error:   true,
				})
			}
			continue
		}

		if msg, ok := checkSchemaType(val, rule.Type); !ok {
			findings = append(findings, SchemaFinding{
				Key:     rule.Key,
				Message: fmt.Sprintf("key %q: %s", rule.Key, msg),
				Error:   true,
			})
		}
	}

	return findings
}

// HasSchemaErrors returns true if any finding is an error.
func HasSchemaErrors(findings []SchemaFinding) bool {
	for _, f := range findings {
		if f.Error {
			return true
		}
	}
	return false
}

func checkSchemaType(val string, t SchemaType) (string, bool) {
	switch t {
	case SchemaTypeInt:
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Sprintf("expected int, got %q", val), false
		}
	case SchemaTypeBool:
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Sprintf("expected bool, got %q", val), false
		}
	case SchemaTypeURL:
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			return fmt.Sprintf("expected URL (http/https), got %q", val), false
		}
	}
	return "", true
}
