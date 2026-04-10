// Package schema provides validation of environment variable values
// against declared type schemas (e.g. int, bool, url, email).
package schema

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Type represents a declared schema type for an env var.
type Type string

const (
	TypeString Type = "string"
	TypeInt    Type = "int"
	TypeBool   Type = "bool"
	TypeURL    Type = "url"
	TypeEmail  Type = "email"
)

// Rule maps an environment variable key to an expected Type.
type Rule struct {
	Key      string
	Type     Type
	Required bool
}

// Finding is a schema validation result for a single key.
type Finding struct {
	Key     string
	Type    Type
	Value   string
	Message string
}

var emailRE = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// Validate checks a map of env vars against a set of schema rules.
// It returns a slice of findings for any violations.
func Validate(env map[string]string, rules []Rule) []Finding {
	var findings []Finding

	for _, rule := range rules {
		val, exists := env[rule.Key]
		if !exists {
			if rule.Required {
				findings = append(findings, Finding{
					Key:     rule.Key,
					Type:    rule.Type,
					Message: "required key is missing",
				})
			}
			continue
		}

		if err := checkType(val, rule.Type); err != nil {
			findings = append(findings, Finding{
				Key:     rule.Key,
				Type:    rule.Type,
				Value:   val,
				Message: err.Error(),
			})
		}
	}

	return findings
}

// HasErrors returns true if any finding is present.
func HasErrors(findings []Finding) bool {
	return len(findings) > 0
}

func checkType(val string, t Type) error {
	switch t {
	case TypeInt:
		if _, err := strconv.Atoi(strings.TrimSpace(val)); err != nil {
			return fmt.Errorf("expected int, got %q", val)
		}
	case TypeBool:
		v := strings.ToLower(strings.TrimSpace(val))
		if v != "true" && v != "false" && v != "1" && v != "0" {
			return fmt.Errorf("expected bool (true/false/1/0), got %q", val)
		}
	case TypeURL:
		u, err := url.ParseRequestURI(val)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("expected valid URL, got %q", val)
		}
	case TypeEmail:
		if !emailRE.MatchString(val) {
			return fmt.Errorf("expected valid email, got %q", val)
		}
	}
	return nil
}
