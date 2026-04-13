package env

import (
	"strings"
	"testing"
)

func TestConvertMap_DotenvFormat(t *testing.T) {
	m := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	out, err := ConvertMap(m, DefaultConvertOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", out)
	}
}

func TestConvertMap_ExportFormat(t *testing.T) {
	m := map[string]string{"DB_HOST": "localhost"}
	opts := DefaultConvertOptions()
	opts.Format = FormatExport
	out, err := ConvertMap(m, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export DB_HOST=localhost") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestConvertMap_JSONFormat(t *testing.T) {
	m := map[string]string{"KEY": "value"}
	opts := DefaultConvertOptions()
	opts.Format = FormatJSON
	out, err := ConvertMap(m, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"KEY": "value"`) {
		t.Errorf("expected JSON key-value pair, got:\n%s", out)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON to start with '{', got:\n%s", out)
	}
}

func TestConvertMap_InlineFormat(t *testing.T) {
	m := map[string]string{"FOO": "bar"}
	opts := DefaultConvertOptions()
	opts.Format = FormatInline
	out, err := ConvertMap(m, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected inline FOO=bar, got:\n%s", out)
	}
	if strings.Count(out, "\n") != 1 {
		t.Errorf("expected single newline for inline format, got:\n%s", out)
	}
}

func TestConvertMap_QuotedValues(t *testing.T) {
	m := map[string]string{"MSG": "hello world"}
	opts := DefaultConvertOptions()
	opts.QuoteValues = true
	out, err := ConvertMap(m, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestConvertMap_UnsupportedFormat_ReturnsError(t *testing.T) {
	m := map[string]string{"X": "1"}
	opts := DefaultConvertOptions()
	opts.Format = ConvertFormat("xml")
	_, err := ConvertMap(m, opts)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestConvertMap_EmptyMap_ReturnsEmptyOrBraces(t *testing.T) {
	m := map[string]string{}
	out, err := ConvertMap(m, DefaultConvertOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty output for empty map, got: %q", out)
	}
}
