package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents a target type for casting.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastBool   CastType = "bool"
)

// CastResult holds the outcome of casting a single key.
type CastResult struct {
	Key      string
	Original string
	Casted   string
	Type     CastType
	Err      error
}

// DefaultCastOptions returns sensible defaults.
func DefaultCastOptions() CastOptions {
	return CastOptions{
		TargetType:   CastString,
		Keys:         nil,
		SkipOnError:  true,
		TrimSpace:    true,
	}
}

// CastOptions configures CastMap behaviour.
type CastOptions struct {
	TargetType  CastType
	Keys        []string // if empty, cast all keys
	SkipOnError bool     // skip keys that fail casting instead of returning error
	TrimSpace   bool
}

// CastMap attempts to cast values in src to the target type.
// It returns a new map with cast values and a slice of results.
func CastMap(src map[string]string, opts CastOptions) (map[string]string, []CastResult, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	var results []CastResult
	for k, v := range src {
		if len(keySet) > 0 {
			if _, ok := keySet[k]; !ok {
				continue
			}
		}
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		casted, err := castValue(v, opts.TargetType)
		r := CastResult{Key: k, Original: src[k], Casted: casted, Type: opts.TargetType, Err: err}
		results = append(results, r)
		if err != nil {
			if !opts.SkipOnError {
				return nil, results, fmt.Errorf("cast %q: %w", k, err)
			}
			continue
		}
		out[k] = casted
	}
	return out, results, nil
}

func castValue(v string, t CastType) (string, error) {
	switch t {
	case CastInt:
		_, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return v, fmt.Errorf("cannot cast %q to int", v)
		}
		return v, nil
	case CastFloat:
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return v, fmt.Errorf("cannot cast %q to float", v)
		}
		return v, nil
	case CastBool:
		norm := strings.ToLower(v)
		switch norm {
		case "true", "1", "yes", "on":
			return "true", nil
		case "false", "0", "no", "off":
			return "false", nil
		default:
			return v, fmt.Errorf("cannot cast %q to bool", v)
		}
	default:
		return v, nil
	}
}

// HasCastErrors returns true if any result contains an error.
func HasCastErrors(results []CastResult) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
