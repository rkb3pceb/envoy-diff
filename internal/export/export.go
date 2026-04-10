// Package export provides functionality to export diff results
// to various file formats such as CSV and Markdown.
package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
)

// Options holds configuration for the export operation.
type Options struct {
	Format    Format
	OldFile   string
	NewFile   string
	Timestamp time.Time
}

// Exporter writes diff results to an output stream.
type Exporter interface {
	Write(w io.Writer, changes []diff.Change, opts Options) error
}

// New returns an Exporter for the given format, or an error if unsupported.
func New(format Format) (Exporter, error) {
	switch format {
	case FormatCSV:
		return &csvExporter{}, nil
	case FormatMarkdown:
		return &markdownExporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// csvExporter exports changes as CSV.
type csvExporter struct{}

func (e *csvExporter) Write(w io.Writer, changes []diff.Change, opts Options) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "type", "old_value", "new_value"}); err != nil {
		return fmt.Errorf("writing csv header: %w", err)
	}
	for _, c := range changes {
		row := []string{c.Key, string(c.Type), c.OldValue, c.NewValue}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("writing csv row for key %q: %w", c.Key, err)
		}
	}
	cw.Flush()
	return cw.Error()
}

// markdownExporter exports changes as a Markdown table.
type markdownExporter struct{}

func (e *markdownExporter) Write(w io.Writer, changes []diff.Change, opts Options) error {
	ts := opts.Timestamp.Format(time.RFC3339)
	fmt.Fprintf(w, "# Env Diff Report\n\n")
	fmt.Fprintf(w, "**Old:** `%s`  \n**New:** `%s`  \n**Generated:** %s\n\n", opts.OldFile, opts.NewFile, ts)

	if len(changes) == 0 {
		fmt.Fprintln(w, "_No changes detected._")
		return nil
	}

	fmt.Fprintln(w, "| Key | Type | Old Value | New Value |")
	fmt.Fprintln(w, "|-----|------|-----------|-----------|")
	for _, c := range changes {
		old := escapeMarkdown(c.OldValue)
		new := escapeMarkdown(c.NewValue)
		fmt.Fprintf(w, "| `%s` | %s | %s | %s |\n", c.Key, string(c.Type), old, new)
	}
	return nil
}

func escapeMarkdown(s string) string {
	if s == "" {
		return "_empty_"
	}
	return strings.ReplaceAll(s, "|", `\|`)
}

// SupportedFormats returns a slice of all supported export format strings.
// This is useful for generating help text or validation messages.
func SupportedFormats() []string {
	return []string{
		string(FormatCSV),
		string(FormatMarkdown),
	}
}
