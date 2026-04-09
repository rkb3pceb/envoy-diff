// Package config handles loading and validating envoy-diff CLI configuration
// from a YAML config file (e.g. ~/.envoy-diff.yaml or .envoy-diff.yaml).
package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds user-defined defaults for envoy-diff behaviour.
type Config struct {
	// DefaultFormat is the output format used when --format is not specified.
	// Accepted values: "text", "json".
	DefaultFormat string `yaml:"default_format"`

	// AuditMode enables audit findings by default when true.
	AuditMode bool `yaml:"audit_mode"`

	// SensitivePatterns is a list of additional key substrings treated as
	// sensitive during auditing (merged with built-in patterns).
	SensitivePatterns []string `yaml:"sensitive_patterns"`

	// HistoryDir overrides the default directory used to store diff history.
	HistoryDir string `yaml:"history_dir"`

	// SnapshotDir overrides the default directory used to store snapshots.
	SnapshotDir string `yaml:"snapshot_dir"`
}

// defaults returns a Config populated with sensible fallback values.
func defaults() Config {
	return Config{
		DefaultFormat: "text",
		AuditMode:     false,
	}
}

// Load reads a Config from the given file path.
// If path is empty, Load searches for .envoy-diff.yaml in the current
// directory and then in the user's home directory, returning defaults
// when no file is found.
func Load(path string) (Config, error) {
	cfg := defaults()

	resolved, err := resolve(path)
	if err != nil {
		// No config file found — return defaults silently.
		return cfg, nil
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	if err := validate(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// resolve returns the first readable config file path.
func resolve(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}

	candidates := []string{".envoy-diff.yaml"}

	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".envoy-diff.yaml"))
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}

	return "", errors.New("no config file found")
}

// validate checks that Config fields hold acceptable values.
func validate(cfg Config) error {
	switch cfg.DefaultFormat {
	case "text", "json", "":
		// valid
	default:
		return errors.New("config: default_format must be \"text\" or \"json\"")
	}
	return nil
}
