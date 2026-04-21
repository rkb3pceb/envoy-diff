package env

import (
	"fmt"
	"strings"
)

// DefaultSplitOptions returns a SplitOptions with safe defaults.
func DefaultSplitOptions() SplitOptions {
	return SplitOptions{
		Delimiter:    ",",
		KeyTemplate:  "{KEY}_{INDEX}",
		IndexBase:    1,
		SkipEmpty:    true,
		TrimParts:    true,
	}
}

// SplitOptions controls how multi-value env vars are split into indexed keys.
type SplitOptions struct {
	// Delimiter is the separator used to split values (default: ",").
	Delimiter string
	// KeyTemplate is the pattern for generated keys; {KEY} and {INDEX} are replaced.
	KeyTemplate string
	// IndexBase is the starting index (0 or 1).
	IndexBase int
	// Keys restricts splitting to these keys only; empty means all keys.
	Keys []string
	// SkipEmpty drops empty parts after splitting.
	SkipEmpty bool
	// TrimParts trims whitespace from each part.
	TrimParts bool
	// KeepOriginal retains the original key alongside the split keys.
	KeepOriginal bool
}

// SplitResult holds the output map and metadata about what was split.
type SplitResult struct {
	Map      map[string]string
	SplitKeys []string
}

// HasSplitChanges returns true when at least one key was split.
func HasSplitChanges(r SplitResult) bool {
	return len(r.SplitKeys) > 0
}

// SplitMap splits multi-value environment variables into indexed keys.
// Keys not targeted by the options are copied verbatim.
func SplitMap(src map[string]string, opts SplitOptions) (SplitResult, error) {
	if opts.Delimiter == "" {
		return SplitResult{}, fmt.Errorf("split: delimiter must not be empty")
	}
	if opts.KeyTemplate == "" {
		opts.KeyTemplate = DefaultSplitOptions().KeyTemplate
	}

	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[k] = true
	}

	out := make(map[string]string, len(src))
	var splitKeys []string

	for k, v := range src {
		if len(targetSet) > 0 && !targetSet[k] {
			out[k] = v
			continue
		}

		parts := strings.Split(v, opts.Delimiter)
		filtered := parts[:0]
		for _, p := range parts {
			if opts.TrimParts {
				p = strings.TrimSpace(p)
			}
			if opts.SkipEmpty && p == "" {
				continue
			}
			filtered = append(filtered, p)
		}

		if len(filtered) <= 1 {
			out[k] = v
			continue
		}

		if opts.KeepOriginal {
			out[k] = v
		}
		splitKeys = append(splitKeys, k)

		for i, part := range filtered {
			newKey := strings.ReplaceAll(opts.KeyTemplate, "{KEY}", k)
			newKey = strings.ReplaceAll(newKey, "{INDEX}", fmt.Sprintf("%d", i+opts.IndexBase))
			out[newKey] = part
		}
	}

	return SplitResult{Map: out, SplitKeys: splitKeys}, nil
}
