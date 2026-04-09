package redact

// PatternList returns a copy of the built-in sensitive key patterns.
// Callers can inspect which patterns are active without being able to
// mutate the package-level slice.
func PatternList() []string {
	copy := make([]string, len(sensitivePatterns))
	for i, p := range sensitivePatterns {
		copy[i] = p
	}
	return copy
}

// Placeholder returns the string used to replace redacted values.
// This is exported so formatters and reporters can reference the same
// constant without importing a magic string.
func Placeholder() string {
	return redactedPlaceholder
}
