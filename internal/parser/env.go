// Package parser provides utilities for parsing environment variable sources.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// EnvMap represents a set of environment variables as key-value pairs.
type EnvMap map[string]string

// ParseEnvFile reads an .env-style file from the given reader and returns
// an EnvMap. Lines starting with '#' are treated as comments and ignored.
// Blank lines are also ignored. Values may optionally be quoted.
func ParseEnvFile(r io.Reader) (EnvMap, error) {
	env := make(EnvMap)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading env file: %w", err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format %q: expected KEY=VALUE", line)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value := strings.TrimSpace(parts[1])
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
