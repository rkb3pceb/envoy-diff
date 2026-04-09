package formatter

import (
	"encoding/json"
	"io"

	"envoy-diff/internal/diff"
)

// JSONFormatter formats diff results as JSON
type JSONFormatter struct{}

// JSONOutput represents the JSON structure for output
type JSONOutput struct {
	Summary Summary  `json:"summary"`
	Changes []Change `json:"changes"`
}

// Summary contains aggregate information
type Summary struct {
	Total    int `json:"total"`
	Added    int `json:"added"`
	Modified int `json:"modified"`
	Removed  int `json:"removed"`
}

// Change represents a single environment variable change
type Change struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// Format writes the diff result in JSON format
func (f *JSONFormatter) Format(result *diff.Result, w io.Writer) error {
	output := f.buildOutput(result)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// buildOutput constructs the JSON output structure
func (f *JSONFormatter) buildOutput(result *diff.Result) JSONOutput {
	output := JSONOutput{
		Summary: Summary{Total: len(result.Changes)},
		Changes: make([]Change, 0, len(result.Changes)),
	}

	for key, changeType := range result.Changes {
		change := Change{
			Key:  key,
			Type: string(changeType),
		}

		switch changeType {
		case diff.Added:
			change.NewValue = result.NewEnv[key]
			output.Summary.Added++
		case diff.Removed:
			change.OldValue = result.OldEnv[key]
			output.Summary.Removed++
		case diff.Modified:
			change.OldValue = result.OldEnv[key]
			change.NewValue = result.NewEnv[key]
			output.Summary.Modified++
		}

		output.Changes = append(output.Changes, change)
	}

	return output
}
