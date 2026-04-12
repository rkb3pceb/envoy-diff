package env

import (
	"fmt"
	"strings"
)

// TransformFn is a function that transforms a key-value pair.
// It returns the new key, new value, and whether to keep the entry.
type TransformFn func(key, value string) (string, string, bool)

// TransformOptions configures the Map transformation pipeline.
type TransformOptions struct {
	// PrefixAdd prepends a string to every key.
	PrefixAdd string

	// PrefixStrip removes a prefix from keys that have it.
	PrefixStrip string

	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool

	// LowercaseKeys converts all keys to lower_case.
	LowercaseKeys bool

	// DropEmpty removes entries with empty values.
	DropEmpty bool

	// Extra holds user-supplied transform functions applied last.
	Extra []TransformFn
}

// DefaultTransformOptions returns a no-op TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// TransformMap applies the configured transformations to src and returns
// a new map. The original map is never mutated.
func TransformMap(src map[string]string, opts TransformOptions) (map[string]string, error) {
	if opts.UppercaseKeys && opts.LowercaseKeys {
		return nil, fmt.Errorf("transform: UppercaseKeys and LowercaseKeys are mutually exclusive")
	}

	out := make(map[string]string, len(src))

	for k, v := range src {
		// strip prefix first so add-prefix applies to the stripped key
		if opts.PrefixStrip != "" {
			k = strings.TrimPrefix(k, opts.PrefixStrip)
		}
		if opts.PrefixAdd != "" {
			k = opts.PrefixAdd + k
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.LowercaseKeys {
			k = strings.ToLower(k)
		}
		if opts.DropEmpty && v == "" {
			continue
		}

		keep := true
		for _, fn := range opts.Extra {
			var nk, nv string
			nk, nv, keep = fn(k, v)
			if !keep {
				break
			}
			k, v = nk, nv
		}
		if keep {
			out[k] = v
		}
	}
	return out, nil
}
