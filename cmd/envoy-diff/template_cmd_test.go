package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplateFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "env.tmpl")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write template file: %v", err)
	}
	return p
}

func TestTemplateCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"template"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no file arg supplied")
	}
}

func TestTemplateCmd_PlainEnvFile(t *testing.T) {
	path := writeTemplateFile(t, "KEY=value\n")
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"template", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY=value") {
		t.Errorf("expected output to contain KEY=value, got: %q", buf.String())
	}
}

func TestTemplateCmd_VarSubstitution(t *testing.T) {
	path := writeTemplateFile(t, "ENV={{ .DEPLOY_ENV }}\n")
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"template", path, "--var", "DEPLOY_ENV=staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ENV=staging") {
		t.Errorf("expected ENV=staging in output, got: %q", buf.String())
	}
}

func TestTemplateCmd_InvalidVar_ReturnsError(t *testing.T) {
	path := writeTemplateFile(t, "X=1\n")
	rootCmd.SetArgs([]string{"template", path, "--var", "NOEQUALS"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --var format")
	}
}

func TestTemplateCmd_OutputFile(t *testing.T) {
	path := writeTemplateFile(t, "APP=envoy\n")
	outPath := filepath.Join(t.TempDir(), "out.env")
	rootCmd.SetArgs([]string{"template", path, "--output", outPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("output file not written: %v", err)
	}
	if !strings.Contains(string(data), "APP=envoy") {
		t.Errorf("unexpected output file content: %q", data)
	}
}
