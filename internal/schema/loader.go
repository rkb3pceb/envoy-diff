package schema

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// fileSchema is the YAML structure for a schema file.
type fileSchema struct {
	Rules []struct {
		Key      string `yaml:"key"`
		Type     string `yaml:"type"`
		Required bool   `yaml:"required"`
	} `yaml:"rules"`
}

// Load reads a YAML schema file and returns a slice of Rules.
// Example schema file:
//
//	rules:
//	  - key: PORT
//	    type: int
//	    required: true
//	  - key: API_URL
//	    type: url
func Load(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: cannot read file %q: %w", path, err)
	}

	var fs fileSchema
	if err := yaml.Unmarshal(data, &fs); err != nil {
		return nil, fmt.Errorf("schema: invalid YAML in %q: %w", path, err)
	}

	var rules []Rule
	for _, r := range fs.Rules {
		if r.Key == "" {
			return nil, fmt.Errorf("schema: rule with empty key in %q", path)
		}
		t, err := parseType(r.Type)
		if err != nil {
			return nil, fmt.Errorf("schema: key %q: %w", r.Key, err)
		}
		rules = append(rules, Rule{
			Key:      r.Key,
			Type:     t,
			Required: r.Required,
		})
	}

	return rules, nil
}

// Empty returns an empty rule set (no-op schema).
func Empty() []Rule { return nil }

func parseType(s string) (Type, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "string":
		return TypeString, nil
	case "int":
		return TypeInt, nil
	case "bool":
		return TypeBool, nil
	case "url":
		return TypeURL, nil
	case "email":
		return TypeEmail, nil
	default:
		return "", fmt.Errorf("unknown schema type %q", s)
	}
}
