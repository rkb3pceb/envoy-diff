package formatter

import (
	"fmt"
	"io"
	"sort"

	"envoy-diff/internal/diff"
)

// TextFormatter formats diff results as human-readable text
type TextFormatter struct{}

// Format writes the diff result in text format
func (f *TextFormatter) Format(result *diff.Result, w io.Writer) error {
	if len(result.Changes) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return nil
	}

	fmt.Fprintf(w, "Environment Variable Changes: %d\n\n", len(result.Changes))

	// Sort changes by key for consistent output
	keys := make([]string, 0, len(result.Changes))
	for key := range result.Changes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Group changes by type
	for _, changeType := range []diff.ChangeType{diff.Added, diff.Modified, diff.Removed} {
		var changes []string
		for _, key := range keys {
			if result.Changes[key] == changeType {
				changes = append(changes, key)
			}
		}

		if len(changes) == 0 {
			continue
		}

		f.formatSection(w, changeType, changes, result)
	}

	return nil
}

// formatSection formats a section of changes
func (f *TextFormatter) formatSection(w io.Writer, ct diff.ChangeType, keys []string, result *diff.Result) {
	title := map[diff.ChangeType]string{
		diff.Added:    "Added Variables",
		diff.Modified: "Modified Variables",
		diff.Removed:  "Removed Variables",
	}[ct]

	fmt.Fprintf(w, "%s (%d):\n", title, len(keys))
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, key := range keys {
		symbol := formatChangeType(ct)
		switch ct {
		case diff.Added:
			value := truncateValue(result.NewEnv[key], 40)
			fmt.Fprintf(w, "%s %s = %s\n", symbol, padRight(key, 25), value)
		case diff.Removed:
			value := truncateValue(result.OldEnv[key], 40)
			fmt.Fprintf(w, "%s %s = %s\n", symbol, padRight(key, 25), value)
		case diff.Modified:
			oldVal := truncateValue(result.OldEnv[key], 35)
			newVal := truncateValue(result.NewEnv[key], 35)
			fmt.Fprintf(w, "%s %s\n", symbol, key)
			fmt.Fprintf(w, "  Old: %s\n", oldVal)
			fmt.Fprintf(w, "  New: %s\n", newVal)
		}
	}
	fmt.Fprintln(w)
}
