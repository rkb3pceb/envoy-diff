package env

import (
	"fmt"
	"sort"
	"strings"
)

// InspectOptions controls how a single key is inspected.
type InspectOptions struct {
	// Redact masks the value if the key is considered sensitive.
	Redact bool
	// IncludeMetadata adds derived metadata fields to the result.
	IncludeMetadata bool
}

// DefaultInspectOptions returns sensible defaults.
func DefaultInspectOptions() InspectOptions {
	return InspectOptions{
		Redact:          true,
		IncludeMetadata: true,
	}
}

// InspectResult holds all information about a single env key.
type InspectResult struct {
	Key         string
	Value       string
	Exists      bool
	Empty       bool
	Sensitive   bool
	Length      int
	Metadata    map[string]string
}

// Inspect returns detailed information about a specific key in the map.
func Inspect(env map[string]string, key string, opts InspectOptions) InspectResult {
	val, exists := env[key]

	sensitive := isSensitiveKey(key)
	displayVal := val
	if exists && sensitive && opts.Redact {
		displayVal = "[REDACTED]"
	}

	result := InspectResult{
		Key:       key,
		Value:     displayVal,
		Exists:    exists,
		Empty:     exists && val == "",
		Sensitive: sensitive,
		Length:    len(val),
		Metadata:  map[string]string{},
	}

	if opts.IncludeMetadata && exists {
		result.Metadata["uppercase"] = fmt.Sprintf("%v", key == strings.ToUpper(key))
		result.Metadata["has_spaces"] = fmt.Sprintf("%v", strings.Contains(val, " "))
		result.Metadata["numeric"] = fmt.Sprintf("%v", isNumeric(val))
		result.Metadata["boolean_like"] = fmt.Sprintf("%v", isBoolLike(val))
	}

	return result
}

// InspectAll returns InspectResult for every key in the map, sorted by key.
func InspectAll(env map[string]string, opts InspectOptions) []InspectResult {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]InspectResult, 0, len(keys))
	for _, k := range keys {
		results = append(results, Inspect(env, k, opts))
	}
	return results
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isBoolLike(s string) bool {
	switch strings.ToLower(s) {
	case "true", "false", "yes", "no", "1", "0", "on", "off":
		return true
	}
	return false
}
