// Package redact provides utilities for masking sensitive values
// before they are displayed or written to output.
package redact

import "strings"

// sensitivePatterns holds substrings that indicate a key is sensitive.
var sensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"ACCESS_KEY",
	"SIGNING_KEY",
}

const redactedPlaceholder = "[REDACTED]"

// IsSensitive reports whether the given key name matches any known
// sensitive pattern (case-insensitive).
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Value returns the original value unchanged if the key is not
// sensitive, or the redacted placeholder if it is.
func Value(key, value string) string {
	if IsSensitive(key) {
		return redactedPlaceholder
	}
	return value
}

// Map returns a copy of the provided map with sensitive values replaced
// by the redacted placeholder.
func Map(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Value(k, v)
	}
	return out
}
