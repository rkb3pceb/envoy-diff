package env

import "fmt"

// CopyOptions controls the behaviour of CopyKeys.
type CopyOptions struct {
	// SourcePrefix filters which keys are copied (empty = all keys).
	SourcePrefix string
	// DestPrefix is prepended to every copied key.
	DestPrefix string
	// Overwrite allows existing destination keys to be replaced.
	Overwrite bool
	// ErrorOnMissing returns an error when a requested key is absent.
	ErrorOnMissing bool
}

// DefaultCopyOptions returns sensible defaults for CopyKeys.
func DefaultCopyOptions() CopyOptions {
	return CopyOptions{
		Overwrite: false,
		ErrorOnMissing: false,
	}
}

// CopyResult records what happened during a copy operation.
type CopyResult struct {
	Copied    []string // keys that were written to dst
	Skipped   []string // keys skipped because dst already had them
	Missing   []string // requested keys not found in src
}

// CopyKeys copies keys from src into dst according to opts.
// When opts.SourcePrefix is set only keys with that prefix are considered.
// The original src map is never mutated.
func CopyKeys(src, dst map[string]string, keys []string, opts CopyOptions) (CopyResult, error) {
	result := CopyResult{}

	if len(keys) == 0 {
		// Copy all keys matching the optional source prefix.
		for k := range src {
			if opts.SourcePrefix == "" || len(k) >= len(opts.SourcePrefix) && k[:len(opts.SourcePrefix)] == opts.SourcePrefix {
				keys = append(keys, k)
			}
		}
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			result.Missing = append(result.Missing, k)
			if opts.ErrorOnMissing {
				return result, fmt.Errorf("copy: key %q not found in source", k)
			}
			continue
		}

		destKey := opts.DestPrefix + k
		if _, exists := dst[destKey]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, destKey)
			continue
		}

		dst[destKey] = v
		result.Copied = append(result.Copied, destKey)
	}

	return result, nil
}

// HasCopied returns true when at least one key was copied.
func HasCopied(r CopyResult) bool { return len(r.Copied) > 0 }
