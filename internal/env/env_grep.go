package env

import (
	"regexp"
	"strings"
)

// GrepOptions controls how GrepMap filters environment variable entries.
type GrepOptions struct {
	// Pattern is a substring or regex pattern to match against keys and/or values.
	Pattern string
	// UseRegex treats Pattern as a regular expression instead of a plain substring.
	UseRegex bool
	// MatchKeys enables matching against variable keys.
	MatchKeys bool
	// MatchValues enables matching against variable values.
	MatchValues bool
	// CaseInsensitive makes the match case-insensitive.
	CaseInsensitive bool
	// Invert returns entries that do NOT match the pattern.
	Invert bool
}

// DefaultGrepOptions returns sensible defaults: case-insensitive substring
// search across both keys and values.
func DefaultGrepOptions() GrepOptions {
	return GrepOptions{
		MatchKeys:       true,
		MatchValues:     true,
		CaseInsensitive: true,
	}
}

// GrepMap filters the provided map, returning only entries whose key or value
// matches the configured pattern. The original map is not mutated.
func GrepMap(src map[string]string, opts GrepOptions) (map[string]string, error) {
	if opts.Pattern == "" {
		// No pattern — return a shallow copy.
		out := make(map[string]string, len(src))
		for k, v := range src {
			out[k] = v
		}
		return out, nil
	}

	matchFn, err := buildMatchFn(opts)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string)
	for k, v := range src {
		hit := false
		if opts.MatchKeys && matchFn(k) {
			hit = true
		}
		if !hit && opts.MatchValues && matchFn(v) {
			hit = true
		}
		if opts.Invert {
			hit = !hit
		}
		if hit {
			out[k] = v
		}
	}
	return out, nil
}

// HasGrepResults returns true when the result map contains at least one entry.
func HasGrepResults(result map[string]string) bool {
	return len(result) > 0
}

func buildMatchFn(opts GrepOptions) (func(string) bool, error) {
	if opts.UseRegex {
		flags := ""
		if opts.CaseInsensitive {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + opts.Pattern)
		if err != nil {
			return nil, err
		}
		return re.MatchString, nil
	}
	pattern := opts.Pattern
	if opts.CaseInsensitive {
		pattern = strings.ToLower(pattern)
		return func(s string) bool {
			return strings.Contains(strings.ToLower(s), pattern)
		}, nil
	}
	return func(s string) bool {
		return strings.Contains(s, pattern)
	}, nil
}
