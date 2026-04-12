package env

import (
	"fmt"
	"strings"
)

// ValidationLevel indicates the severity of a validation finding.
type ValidationLevel string

const (
	LevelError   ValidationLevel = "error"
	LevelWarning ValidationLevel = "warning"
)

// ValidationFinding represents a single issue found during env map validation.
type ValidationFinding struct {
	Key     string
	Message string
	Level   ValidationLevel
}

// ValidateOptions controls which checks are applied during validation.
type ValidateOptions struct {
	// RequireUppercase causes a warning when a key contains lowercase letters.
	RequireUppercase bool
	// ForbidEmpty causes an error when a value is empty.
	ForbidEmpty bool
	// MaxKeyLength, when > 0, causes an error if a key exceeds this length.
	MaxKeyLength int
	// AllowedPrefixes, when non-empty, warns on keys that match none of them.
	AllowedPrefixes []string
}

// DefaultValidateOptions returns a ValidateOptions with sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		RequireUppercase: true,
		ForbidEmpty:      false,
		MaxKeyLength:     128,
	}
}

// ValidateMap checks each key/value pair in m against opts and returns
// all findings. An empty slice means no issues were found.
func ValidateMap(m map[string]string, opts ValidateOptions) []ValidationFinding {
	var findings []ValidationFinding

	for k, v := range m {
		if opts.RequireUppercase && k != strings.ToUpper(k) {
			findings = append(findings, ValidationFinding{
				Key:     k,
				Message: "key contains lowercase letters",
				Level:   LevelWarning,
			})
		}

		if opts.ForbidEmpty && v == "" {
			findings = append(findings, ValidationFinding{
				Key:     k,
				Message: "value is empty",
				Level:   LevelError,
			})
		}

		if opts.MaxKeyLength > 0 && len(k) > opts.MaxKeyLength {
			findings = append(findings, ValidationFinding{
				Key:     k,
				Message: fmt.Sprintf("key length %d exceeds maximum %d", len(k), opts.MaxKeyLength),
				Level:   LevelError,
			})
		}

		if len(opts.AllowedPrefixes) > 0 && !hasAllowedPrefix(k, opts.AllowedPrefixes) {
			findings = append(findings, ValidationFinding{
				Key:     k,
				Message: fmt.Sprintf("key does not match any allowed prefix: %s", strings.Join(opts.AllowedPrefixes, ", ")),
				Level:   LevelWarning,
			})
		}
	}

	return findings
}

// HasErrors returns true if any finding has level error.
func HasValidationErrors(findings []ValidationFinding) bool {
	for _, f := range findings {
		if f.Level == LevelError {
			return true
		}
	}
	return false
}

func hasAllowedPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
