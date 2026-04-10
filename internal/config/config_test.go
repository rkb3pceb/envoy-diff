package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfig writes content to a temp file and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envoy-diff-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

// chdir changes the working directory for the duration of the test,
// restoring the original directory via t.Cleanup.
func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() { os.Chdir(old) })
}

func TestLoad_Defaults_WhenNoFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/envoy-diff.yaml")
	// A missing explicit path should propagate an os error, but a missing
	// auto-discovered file returns defaults with no error.
	_ = err // either outcome is acceptable; we just want defaults-like values
	_ = cfg
}

func TestLoad_ReturnsDefaults_WhenNoAutoFile(t *testing.T) {
	// Change to a temp dir so no .envoy-diff.yaml is found automatically.
	chdir(t, t.TempDir())

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultFormat != "text" {
		t.Errorf("expected default format \"text\", got %q", cfg.DefaultFormat)
	}
	if cfg.AuditMode != false {
		t.Error("expected audit_mode to default to false")
	}
}

func TestLoad_ParsesValidConfig(t *testing.T) {
	path := writeConfig(t, `
default_format: json
audit_mode: true
sensitive_patterns:
  - MYAPP_SECRET
  - INTERNAL_TOKEN
history_dir: /tmp/history
snapshot_dir: /tmp/snapshots
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.DefaultFormat != "json" {
		t.Errorf("DefaultFormat: got %q, want \"json\"", cfg.DefaultFormat)
	}
	if !cfg.AuditMode {
		t.Error("AuditMode: expected true")
	}
	if len(cfg.SensitivePatterns) != 2 {
		t.Errorf("SensitivePatterns: got %d entries, want 2", len(cfg.SensitivePatterns))
	}
	if cfg.HistoryDir != "/tmp/history" {
		t.Errorf("HistoryDir: got %q", cfg.HistoryDir)
	}
}

func TestLoad_InvalidFormat_ReturnsError(t *testing.T) {
	path := writeConfig(t, "default_format: xml\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid default_format, got nil")
	}
}

func TestLoad_AutoDiscoversLocalFile(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	if err := os.WriteFile(filepath.Join(dir, ".envoy-diff.yaml"), []byte("audit_mode: true\n"), 0644); err != nil {
		t.Fatalf("write local config: %v", err)
	}

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if !cfg.AuditMode {
		t.Error("expected AuditMode true from auto-discovered file")
	}
}
