package env

import (
	"fmt"
	"strings"
	"unicode"
)

// LintSeverity indicates the severity of a lint finding.
type LintSeverity string

const (
	LintError   LintSeverity = "error"
	LintWarning LintSeverity = "warning"
)

// LintFinding represents a single lint result for a key/value pair.
type LintFinding struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintOptions controls which lint checks are applied.
type LintOptions struct {
	ForbidEmpty        bool
	ForbidLowercase    bool
	ForbidWhitespace   bool
	ForbidPlaceholders bool
	MaxKeyLength       int
}

// DefaultLintOptions returns sensible defaults for linting.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		ForbidEmpty:        true,
		ForbidLowercase:    true,
		ForbidWhitespace:   true,
		ForbidPlaceholders: true,
		MaxKeyLength:       64,
	}
}

// placeholders are common stub values that indicate unset configuration.
var placeholders = []string{"todo", "changeme", "fixme", "placeholder", "xxx", "tbd"}

// LintMap runs lint checks against a flat env map and returns all findings.
func LintMap(env map[string]string, opts LintOptions) []LintFinding {
	var findings []LintFinding

	for k, v := range env {
		if opts.ForbidWhitespace && strings.ContainsAny(k, " \t") {
			findings = append(findings, LintFinding{
				Key:      k,
				Message:  "key contains whitespace",
				Severity: LintError,
			})
		}

		if opts.ForbidLowercase && hasLowercase(k) {
			findings = append(findings, LintFinding{
				Key:      k,
				Message:  "key contains lowercase letters",
				Severity: LintWarning,
			})
		}

		if opts.MaxKeyLength > 0 && len(k) > opts.MaxKeyLength {
			findings = append(findings, LintFinding{
				Key:      k,
				Message:  fmt.Sprintf("key length %d exceeds maximum %d", len(k), opts.MaxKeyLength),
				Severity: LintError,
			})
		}

		if opts.ForbidEmpty && strings.TrimSpace(v) == "" {
			findings = append(findings, LintFinding{
				Key:      k,
				Message:  "value is empty",
				Severity: LintWarning,
			})
		}

		if opts.ForbidPlaceholders && isPlaceholder(v) {
			findings = append(findings, LintFinding{
				Key:      k,
				Message:  fmt.Sprintf("value looks like a placeholder: %q", v),
				Severity: LintWarning,
			})
		}
	}

	return findings
}

// HasLintErrors returns true if any finding has error severity.
func HasLintErrors(findings []LintFinding) bool {
	for _, f := range findings {
		if f.Severity == LintError {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func isPlaceholder(v string) bool {
	lower := strings.ToLower(strings.TrimSpace(v))
	for _, p := range placeholders {
		if lower == p {
			return true
		}
	}
	return false
}
