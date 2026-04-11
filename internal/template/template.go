// Package template provides env file rendering from Go templates.
// It supports variable substitution and conditional blocks.
package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	gotemplate "text/template"
)

// Options controls template rendering behaviour.
type Options struct {
	// Vars are additional key-value pairs injected into the template context.
	Vars map[string]string
	// AllowMissing suppresses errors for undefined template variables.
	AllowMissing bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Vars:         map[string]string{},
		AllowMissing: false,
	}
}

// Render parses src as a Go template and executes it with the supplied
// options merged with the current OS environment.
func Render(src []byte, opts Options) ([]byte, error) {
	data := buildContext(opts)

	missingKey := "error"
	if opts.AllowMissing {
		missingKey = "zero"
	}

	tmpl, err := gotemplate.New("env").
		Option("missingkey=" + missingKey).
		Funcs(funcMap()).
		Parse(string(src))
	if err != nil {
		return nil, fmt.Errorf("template parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("template execute: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderFile reads path, renders it and returns the result.
func RenderFile(path string, opts Options) ([]byte, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template file: %w", err)
	}
	return Render(src, opts)
}

// buildContext merges OS env with user-supplied vars (user vars win).
func buildContext(opts Options) map[string]string {
	ctx := make(map[string]string)
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			ctx[parts[0]] = parts[1]
		}
	}
	for k, v := range opts.Vars {
		ctx[k] = v
	}
	return ctx
}

// funcMap returns template helper functions.
func funcMap() gotemplate.FuncMap {
	return gotemplate.FuncMap{
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
}
