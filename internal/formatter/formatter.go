// Package formatter provides output formatting for diff results.
package formatter

import (
	"fmt"
	"io"
	"strings"

	"envoy-diff/internal/diff"
)

// Format defines the output format type
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter handles formatting of diff results
type Formatter interface {
	Format(result *diff.Result, w io.Writer) error
}

// New creates a new formatter based on the specified format
func New(format Format) (Formatter, error) {
	switch format {
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}

// colorize adds ANSI color codes to text
func colorize(text, color string) string {
	colors := map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"reset":  "\033[0m",
	}
	if code, ok := colors[color]; ok {
		return code + text + colors["reset"]
	}
	return text
}

// formatChangeType returns a colored symbol for the change type
func formatChangeType(ct diff.ChangeType) string {
	switch ct {
	case diff.Added:
		return colorize("+", "green")
	case diff.Removed:
		return colorize("-", "red")
	case diff.Modified:
		return colorize("~", "yellow")
	default:
		return " "
	}
}

// truncateValue truncates long values for display
func truncateValue(value string, maxLen int) string {
	if len(value) <= maxLen {
		return value
	}
	return value[:maxLen-3] + "..."
}

// padRight pads a string to the right with spaces
func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}
