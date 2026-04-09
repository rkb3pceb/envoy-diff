package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"envoy-diff/internal/diff"
)

func TestNew_ValidFormats(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{"text format", FormatText, "*formatter.TextFormatter"},
		{"json format", FormatJSON, "*formatter.JSONFormatter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.format)
			if err != nil {
				t.Errorf("New() error = %v", err)
			}
			if f == nil {
				t.Error("New() returned nil formatter")
			}
		})
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("invalid")
	if err == nil {
		t.Error("New() expected error for invalid format")
	}
}

func TestTextFormatter_NoChanges(t *testing.T) {
	result := &diff.Result{
		Changes: map[string]diff.ChangeType{},
	}

	var buf bytes.Buffer
	f := &TextFormatter{}
	err := f.Format(result, &buf)

	if err != nil {
		t.Errorf("Format() error = %v", err)
	}

	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("Expected 'No changes detected', got: %s", buf.String())
	}
}

func TestTextFormatter_WithChanges(t *testing.T) {
	result := &diff.Result{
		OldEnv: map[string]string{"DB_HOST": "localhost"},
		NewEnv: map[string]string{"DB_HOST": "prod-db", "API_KEY": "secret"},
		Changes: map[string]diff.ChangeType{
			"DB_HOST": diff.Modified,
			"API_KEY": diff.Added,
		},
	}

	var buf bytes.Buffer
	f := &TextFormatter{}
	err := f.Format(result, &buf)

	if err != nil {
		t.Errorf("Format() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "DB_HOST") || !strings.Contains(output, "API_KEY") {
		t.Errorf("Output missing expected keys: %s", output)
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	result := &diff.Result{
		OldEnv: map[string]string{"KEY1": "old"},
		NewEnv: map[string]string{"KEY1": "new", "KEY2": "added"},
		Changes: map[string]diff.ChangeType{
			"KEY1": diff.Modified,
			"KEY2": diff.Added,
		},
	}

	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.Format(result, &buf)

	if err != nil {
		t.Errorf("Format() error = %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	if output.Summary.Total != 2 {
		t.Errorf("Expected 2 total changes, got %d", output.Summary.Total)
	}

	if output.Summary.Added != 1 || output.Summary.Modified != 1 {
		t.Errorf("Expected 1 added and 1 modified, got %d added, %d modified",
			output.Summary.Added, output.Summary.Modified)
	}
}
