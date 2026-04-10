// Package redact provides utilities for identifying and masking sensitive
// environment variable values before display or export.
package redact

import (
	"strings"
)

// MaskLevel controls how aggressively values are masked.
type MaskLevel int

const (
	// MaskFull replaces the entire value with the placeholder.
	MaskFull MaskLevel = iota
	// MaskPartial reveals the first and last two characters of the value.
	MaskPartial
)

// Mask returns a masked representation of value according to the given level.
// If the value is empty or shorter than 6 characters, MaskFull is always used.
func Mask(value string, level MaskLevel) string {
	if value == "" {
		return Placeholder
	}
	if level == MaskPartial && len(value) >= 6 {
		return string(value[0:2]) + strings.Repeat("*", len(value)-4) + string(value[len(value)-2:])
	}
	return Placeholder
}

// MaskMap returns a copy of m where all sensitive keys have their values
// masked according to the given level.
func MaskMap(m map[string]string, level MaskLevel) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if IsSensitive(k) {
			out[k] = Mask(v, level)
		} else {
			out[k] = v
		}
	}
	return out
}
