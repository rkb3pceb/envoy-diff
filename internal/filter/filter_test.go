package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

var sampleEntries = []diff.Entry{
	{Key: "APP_HOST", OldValue: "localhost", NewValue: "prod.example.com", Type: diff.Modified},
	{Key: "APP_PORT", OldValue: "8080", NewValue: "8080", Type: diff.Unchanged},
	{Key: "DB_PASSWORD", OldValue: "", NewValue: "secret", Type: diff.Added},
	{Key: "LEGACY_FLAG", OldValue: "true", NewValue: "", Type: diff.Removed},
	{Key: "LOG_LEVEL", OldValue: "debug", NewValue: "info", Type: diff.Modified},
}

func TestApply_NoOptions(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{})
	if len(result) != len(sampleEntries) {
		t.Errorf("expected %d entries, got %d", len(sampleEntries), len(result))
	}
}

func TestApply_OnlyChanged(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{OnlyChanged: true})
	for _, e := range result {
		if e.Type == diff.Unchanged {
			t.Errorf("unexpected Unchanged entry: %s", e.Key)
		}
	}
	if len(result) != 4 {
		t.Errorf("expected 4 changed entries, got %d", len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{Prefix: "APP_"})
	if len(result) != 2 {
		t.Errorf("expected 2 entries with prefix APP_, got %d", len(result))
	}
}

func TestApply_KeyContains(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{KeyContains: "LOG"})
	if len(result) != 1 || result[0].Key != "LOG_LEVEL" {
		t.Errorf("expected only LOG_LEVEL, got %+v", result)
	}
}

func TestApply_TypeFilter(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{
		Types: []diff.ChangeType{diff.Added, diff.Removed},
	})
	if len(result) != 2 {
		t.Errorf("expected 2 entries (Added+Removed), got %d", len(result))
	}
	for _, e := range result {
		if e.Type != diff.Added && e.Type != diff.Removed {
			t.Errorf("unexpected type %v for key %s", e.Type, e.Key)
		}
	}
}

func TestApply_CombinedOptions(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{
		Prefix:      "APP_",
		OnlyChanged: true,
	})
	if len(result) != 1 || result[0].Key != "APP_HOST" {
		t.Errorf("expected only APP_HOST, got %+v", result)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	result := filter.Apply(nil, filter.Options{OnlyChanged: true, Prefix: "X_"})
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %+v", result)
	}
}
