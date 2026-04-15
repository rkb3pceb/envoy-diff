package env

import (
	"testing"
)

func TestGrepMap_EmptyPattern_ReturnsAll(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "DB_PASS": "secret"}
	opts := DefaultGrepOptions()
	opts.Pattern = ""
	out, err := GrepMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestGrepMap_SubstringMatch_Keys(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "DB_PASS": "secret", "APP_PORT": "8080"}
	opts := DefaultGrepOptions()
	opts.Pattern = "APP"
	opts.MatchValues = false
	out, err := GrepMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should not be in result")
	}
}

func TestGrepMap_SubstringMatch_Values(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "URL": "https://example.com", "PORT": "9090"}
	opts := DefaultGrepOptions()
	opts.Pattern = "local"
	opts.MatchKeys = false
	out, err := GrepMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost")
	}
}

func TestGrepMap_CaseInsensitive(t *testing.T) {
	src := map[string]string{"APP_ENV": "Production", "LOG_LEVEL": "debug"}
	opts := DefaultGrepOptions()
	opts.Pattern = "production"
	opts.MatchKeys = false
	out, err := GrepMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestGrepMap_RegexMatch(t *testing.T) {
	src := map[string]string{"DB_HOST": "db1.internal", "DB_PORT": "5432", "APP_NAME": "envoy"}
	opts := DefaultGrepOptions()
	opts.Pattern = `^DB_`
	opts.UseRegex = true
	opts.MatchValues = false
	out, err := GrepMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestGrepMap_InvalidRegex_ReturnsError(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	opts := DefaultGrepOptions()
	opts.Pattern = `[invalid`
	opts.UseRegex = true
	_, err := GrepMap(src, opts)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestGrepMap_Invert_ReturnsNonMatching(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "DB_PASS": "secret", "APP_PORT": "8080"}
	opts := DefaultGrepOptions()
	opts.Pattern = "APP"
	opts.MatchValues = false
	opts.Invert = true
	out, err := GrepMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; !ok {
		t.Error("DB_PASS should be in inverted result")
	}
}

func TestHasGrepResults_TrueWhenNonEmpty(t *testing.T) {
	if !HasGrepResults(map[string]string{"K": "V"}) {
		t.Error("expected true for non-empty map")
	}
}

func TestHasGrepResults_FalseWhenEmpty(t *testing.T) {
	if HasGrepResults(map[string]string{}) {
		t.Error("expected false for empty map")
	}
}
