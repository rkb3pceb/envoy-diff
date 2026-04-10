// Package validate provides checks to ensure environment variable maps
// conform to required key presence and value constraints before deployment.
package validate

import "fmt"

// Severity indicates how critical a validation finding is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Finding represents a single validation result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Rule defines a validation rule applied to an env map.
type Rule struct {
	// Required keys that must be present and non-empty.
	RequiredKeys []string
	// ForbiddenKeys that must not appear in the env map.
	ForbiddenKeys []string
	// AllowedPrefixes restricts keys to a set of allowed prefixes.
	// If empty, all prefixes are allowed.
	AllowedPrefixes []string
}

// Validate applies r to the given env map and returns all findings.
func Validate(env map[string]string, r Rule) []Finding {
	var findings []Finding

	for _, key := range r.RequiredKeys {
		v, ok := env[key]
		if !ok || v == "" {
			findings = append(findings, Finding{
				Key:      key,
				Message:  fmt.Sprintf("required key %q is missing or empty", key),
				Severity: SeverityError,
			})
		}
	}

	for _, key := range r.ForbiddenKeys {
		if _, ok := env[key]; ok {
			findings = append(findings, Finding{
				Key:      key,
				Message:  fmt.Sprintf("forbidden key %q must not be present", key),
				Severity: SeverityError,
			})
		}
	}

	if len(r.AllowedPrefixes) > 0 {
		for key := range env {
			if !hasAllowedPrefix(key, r.AllowedPrefixes) {
				findings = append(findings, Finding{
					Key:     key,
					Message: fmt.Sprintf("key %q does not match any allowed prefix", key),
					Severity: SeverityWarning,
				})
			}
		}
	}

	return findings
}

// HasErrors returns true if any finding has error severity.
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}

func hasAllowedPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
