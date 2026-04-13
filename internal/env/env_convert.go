package env

import (
	"fmt"
	"strings"
)

// ConvertFormat represents a supported output serialization format.
type ConvertFormat string

const (
	FormatDotenv  ConvertFormat = "dotenv"
	FormatExport  ConvertFormat = "export"
	FormatJSON    ConvertFormat = "json"
	FormatInline  ConvertFormat = "inline"
)

// DefaultConvertOptions returns options with dotenv as the default output format.
func DefaultConvertOptions() ConvertOptions {
	return ConvertOptions{
		Format: FormatDotenv,
		QuoteValues: false,
		Sorted: true,
	}
}

// ConvertOptions controls how the map is serialized.
type ConvertOptions struct {
	Format      ConvertFormat
	QuoteValues bool
	Sorted      bool
}

// ConvertMap serializes an env map into the requested text format.
// Returns the formatted string and any error encountered.
func ConvertMap(m map[string]string, opts ConvertOptions) (string, error) {
	keys := SortedKeys(m, DefaultSortOptions())
	if !opts.Sorted {
		keys = unsortedKeys(m)
	}

	var sb strings.Builder

	switch opts.Format {
	case FormatDotenv:
		for _, k := range keys {
			v := m[k]
			if opts.QuoteValues {
				v = fmt.Sprintf("%q", v)
			}
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	case FormatExport:
		for _, k := range keys {
			v := m[k]
			if opts.QuoteValues {
				v = fmt.Sprintf("%q", v)
			}
			fmt.Fprintf(&sb, "export %s=%s\n", k, v)
		}
	case FormatJSON:
		sb.WriteString("{\n")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", k, m[k], comma)
		}
		sb.WriteString("}\n")
	case FormatInline:
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			v := m[k]
			if opts.QuoteValues {
				v = fmt.Sprintf("%q", v)
			}
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		sb.WriteString(strings.Join(parts, " "))
		sb.WriteString("\n")
	default:
		return "", fmt.Errorf("unsupported format: %q", opts.Format)
	}

	return sb.String(), nil
}

// unsortedKeys returns map keys in iteration order (non-deterministic).
func unsortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
