package env

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncodeFormat specifies the encoding format to apply to values.
type EncodeFormat string

const (
	EncodeBase64    EncodeFormat = "base64"
	EncodeBase64URL EncodeFormat = "base64url"
	EncodeHex       EncodeFormat = "hex"
)

// DefaultEncodeOptions returns sensible defaults for EncodeMap.
func DefaultEncodeOptions() EncodeOptions {
	return EncodeOptions{
		Format:   EncodeBase64,
		Keys:     nil,
		Decode:   false,
		SkipEmpty: true,
	}
}

// EncodeOptions controls how EncodeMap behaves.
type EncodeOptions struct {
	// Format is the encoding/decoding format.
	Format EncodeFormat
	// Keys restricts encoding to the specified keys; empty means all keys.
	Keys []string
	// Decode reverses the operation (decode instead of encode).
	Decode bool
	// SkipEmpty leaves empty values untouched.
	SkipEmpty bool
}

// EncodeResult holds the output of EncodeMap.
type EncodeResult struct {
	Map      map[string]string
	Encoded  []string
	Errors   []string
}

// EncodeMap encodes or decodes values in env according to opts.
func EncodeMap(env map[string]string, opts EncodeOptions) EncodeResult {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	targetKeys := buildEncodeKeySet(opts.Keys)
	var encoded []string
	var errors []string

	for k, v := range env {
		if len(targetKeys) > 0 && !targetKeys[k] {
			continue
		}
		if opts.SkipEmpty && v == "" {
			continue
		}
		var result string
		var err error
		if opts.Decode {
			result, err = decodeValue(v, opts.Format)
		} else {
			result, err = encodeValue(v, opts.Format)
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", k, err))
			continue
		}
		out[k] = result
		encoded = append(encoded, k)
	}
	return EncodeResult{Map: out, Encoded: encoded, Errors: errors}
}

func encodeValue(v string, format EncodeFormat) (string, error) {
	switch format {
	case EncodeBase64:
		return base64.StdEncoding.EncodeToString([]byte(v)), nil
	case EncodeBase64URL:
		return base64.URLEncoding.EncodeToString([]byte(v)), nil
	case EncodeHex:
		return fmt.Sprintf("%x", []byte(v)), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

func decodeValue(v string, format EncodeFormat) (string, error) {
	switch format {
	case EncodeBase64:
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case EncodeBase64URL:
		b, err := base64.URLEncoding.DecodeString(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case EncodeHex:
		v = strings.TrimPrefix(v, "0x")
		var decoded []byte
		for i := 0; i+1 < len(v); i += 2 {
			var b byte
			_, err := fmt.Sscanf(v[i:i+2], "%02x", &b)
			if err != nil {
				return "", fmt.Errorf("invalid hex at position %d", i)
			}
			decoded = append(decoded, b)
		}
		return string(decoded), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

func buildEncodeKeySet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
