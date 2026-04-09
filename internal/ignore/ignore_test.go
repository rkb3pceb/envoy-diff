package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/ignore"
)

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envoyignore")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write ignore file: %v", err)
	}
	return p
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := ignore.Load("/nonexistent/.envoyignore")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestMatch_ExactKey(t *testing.T) {
	p := writeIgnoreFile(t, "SECRET_KEY\nDEBUG\n")
	r, err := ignore.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !r.Match("SECRET_KEY") {
		t.Error("expected SECRET_KEY to match")
	}
	if !r.Match("DEBUG") {
		t.Error("expected DEBUG to match")
	}
	if r.Match("DATABASE_URL") {
		t.Error("expected DATABASE_URL not to match")
	}
}

func TestMatch_PrefixWildcard(t *testing.T) {
	p := writeIgnoreFile(t, "CI_*\n")
	r, err := ignore.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !r.Match("CI_TOKEN") {
		t.Error("expected CI_TOKEN to match prefix CI_*")
	}
	if !r.Match("CI_BUILD_ID") {
		t.Error("expected CI_BUILD_ID to match prefix CI_*")
	}
	if r.Match("DATABASE_URL") {
		t.Error("expected DATABASE_URL not to match")
	}
}

func TestMatch_IgnoresCommentsAndBlankLines(t *testing.T) {
	p := writeIgnoreFile(t, "# this is a comment\n\nSECRET_KEY\n")
	r, err := ignore.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if r.Match("#") {
		t.Error("comment line should not be treated as a key")
	}
	if !r.Match("SECRET_KEY") {
		t.Error("expected SECRET_KEY to match")
	}
}

func TestEmpty_MatchesNothing(t *testing.T) {
	r := ignore.Empty()
	if r.Match("ANYTHING") {
		t.Error("empty rules should not match any key")
	}
}

func TestFilterKeys_RemovesMatchedKeys(t *testing.T) {
	p := writeIgnoreFile(t, "CI_*\nDEBUG\n")
	r, err := ignore.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	input := []string{"APP_ENV", "CI_TOKEN", "DEBUG", "DATABASE_URL"}
	got := r.FilterKeys(input)
	expected := []string{"APP_ENV", "DATABASE_URL"}
	if len(got) != len(expected) {
		t.Fatalf("FilterKeys returned %v, want %v", got, expected)
	}
	for i, k := range got {
		if k != expected[i] {
			t.Errorf("FilterKeys[%d] = %q, want %q", i, k, expected[i])
		}
	}
}
