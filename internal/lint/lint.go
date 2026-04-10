// Package lint provides rule-based linting for environment variable keys and values.
// It checks for common issues such as lowercase keys, empty values on required vars,
// and suspicious placeholder values that may indicate misconfiguration.
package lint

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Finding represents a single lint issue found in an env map.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// placeholders are values that suggest a variable was never properly set.
var placeholders = []string{
	"TODO", "FIXME", "CHANGEME", "YOUR_", "<", ">", "example", "placeholder",
}

// Lint runs all lint rules against the provided env map and returns any findings.
func Lint(env map[string]string) []Finding {
	var findings []Finding

	for k, v := range env {
		findings = append(findings, checkKey(k, v)...)
		findings = append(findings, checkValue(k, v)...)
	}

	return findings
}

func checkKey(k, _ string) []Finding {
	var findings []Finding

	if k != strings.ToUpper(k) {
		findings = append(findings, Finding{
			Key:      k,
			Message:  "key is not uppercase; convention requires UPPER_SNAKE_CASE",
			Severity: SeverityWarn,
		})
	}

	if strings.Contains(k, " ") {
		findings = append(findings, Finding{
			Key:      k,
			Message:  "key contains whitespace",
			Severity: SeverityError,
		})
	}

	return findings
}

func checkValue(k, v string) []Finding {
	var findings []Finding

	if v == "" {
		findings = append(findings, Finding{
			Key:      k,
			Message:  "value is empty",
			Severity: SeverityWarn,
		})
		return findings
	}

	upper := strings.ToUpper(v)
	for _, p := range placeholders {
		if strings.Contains(upper, strings.ToUpper(p)) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  fmt.Sprintf("value looks like a placeholder (contains %q)", p),
				Severity: SeverityWarn,
			})
			break
		}
	}

	return findings
}
