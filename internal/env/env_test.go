package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envoy-diff/internal/env"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestLoadSources_ParsesFiles(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, "a.env", "FOO=bar\nBAZ=qux\n")

	srcs, err := env.LoadSources([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(srcs) != 1 {
		t.Fatalf("expected 1 source, got %d", len(srcs))
	}
	if srcs[0].Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", srcs[0].Vars["FOO"])
	}
}

func TestLoadSources_MissingFile_ReturnsError(t *testing.T) {
	_, err := env.LoadSources([]string{"/nonexistent/path.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFlatten_LastWins(t *testing.T) {
	srcs := []env.Source{
		{Name: "base", Vars: map[string]string{"FOO": "base", "BAR": "1"}},
		{Name: "override", Vars: map[string]string{"FOO": "override"}},
	}
	result, overridden := env.Flatten(srcs)

	if result["FOO"] != "override" {
		t.Errorf("expected FOO=override, got %q", result["FOO"])
	}
	if result["BAR"] != "1" {
		t.Errorf("expected BAR=1, got %q", result["BAR"])
	}
	if len(overridden) != 1 || overridden[0] != "FOO" {
		t.Errorf("expected [FOO] overridden, got %v", overridden)
	}
}

func TestFlatten_NoConflicts(t *testing.T) {
	srcs := []env.Source{
		{Name: "a", Vars: map[string]string{"A": "1"}},
		{Name: "b", Vars: map[string]string{"B": "2"}},
	}
	result, overridden := env.Flatten(srcs)

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if len(overridden) != 0 {
		t.Errorf("expected no overrides, got %v", overridden)
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	vars := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	keys := env.Keys(vars)

	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
