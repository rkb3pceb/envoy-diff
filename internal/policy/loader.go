package policy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// File represents the YAML structure of a policy file.
type File struct {
	Version string `yaml:"version"`
	Rules   []Rule `yaml:"rules"`
}

// Load reads and parses a policy YAML file from the given path.
func Load(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read %q: %w", path, err)
	}
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("policy: parse %q: %w", path, err)
	}
	if err := validateFile(f); err != nil {
		return nil, fmt.Errorf("policy: invalid %q: %w", path, err)
	}
	return &Policy{Rules: f.Rules}, nil
}

// Empty returns a Policy with no rules (permissive by default).
func Empty() *Policy {
	return &Policy{}
}

func validateFile(f File) error {
	for i, r := range f.Rules {
		if r.Name == "" {
			return fmt.Errorf("rule[%d]: name is required", i)
		}
		if len(r.Keys) == 0 {
			return fmt.Errorf("rule %q: at least one key is required", r.Name)
		}
		if r.Severity != SeverityBlock && r.Severity != SeverityWarn {
			return fmt.Errorf("rule %q: severity must be 'block' or 'warn', got %q", r.Name, r.Severity)
		}
		for _, t := range r.Types {
			switch t {
			case "added", "removed", "modified":
			default:
				return fmt.Errorf("rule %q: unknown type %q", r.Name, t)
			}
		}
	}
	return nil
}
