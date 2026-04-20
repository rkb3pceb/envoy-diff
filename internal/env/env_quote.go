package env

import (
	"fmt"
	"strings"
)

// QuoteStyle controls how values are quoted in output.
type QuoteStyle string

const (
	QuoteStyleNone   QuoteStyle = "none"
	QuoteStyleSingle QuoteStyle = "single"
	QuoteStyleDouble QuoteStyle = "double"
	QuoteStyleAuto   QuoteStyle = "auto" // quote only if value contains spaces or special chars
)

// QuoteOptions configures the quoting behaviour.
type QuoteOptions struct {
	Style      QuoteStyle
	Keys       []string // if non-empty, only quote these keys
	SkipEmpty  bool     // do not quote empty values
}

// DefaultQuoteOptions returns a QuoteOptions with sensible defaults.
func DefaultQuoteOptions() QuoteOptions {
	return QuoteOptions{
		Style:     QuoteStyleAuto,
		SkipEmpty: true,
	}
}

// QuoteMap returns a copy of m with values quoted according to opts.
func QuoteMap(m map[string]string, opts QuoteOptions) (map[string]string, error) {
	if opts.Style == "" {
		opts.Style = QuoteStyleAuto
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if len(keySet) > 0 {
			if _, ok := keySet[k]; !ok {
				out[k] = v
				continue
			}
		}
		if opts.SkipEmpty && v == "" {
			out[k] = v
			continue
		}
		quoted, err := applyQuoteStyle(v, opts.Style)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = quoted
	}
	return out, nil
}

// HasQuoteChanges returns true if any value in quoted differs from original.
func HasQuoteChanges(original, quoted map[string]string) bool {
	for k, v := range quoted {
		if original[k] != v {
			return true
		}
	}
	return false
}

func applyQuoteStyle(v string, style QuoteStyle) (string, error) {
	switch style {
	case QuoteStyleNone:
		return unquote(v), nil
	case QuoteStyleSingle:
		return "'" + strings.ReplaceAll(unquote(v), "'", `'\''`) + "'", nil
	case QuoteStyleDouble:
		return `"` + strings.ReplaceAll(unquote(v), `"`, `\"`) + `"`, nil
	case QuoteStyleAuto:
		raw := unquote(v)
		if needsQuoting(raw) {
			return `"` + strings.ReplaceAll(raw, `"`, `\"`) + `"`, nil
		}
		return raw, nil
	default:
		return "", fmt.Errorf("unknown quote style %q", style)
	}
}

func needsQuoting(v string) bool {
	for _, ch := range v {
		if ch == ' ' || ch == '\t' || ch == '#' || ch == '$' || ch == '\\' || ch == '"' || ch == '\'' {
			return true
		}
	}
	return false
}

func unquote(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}
