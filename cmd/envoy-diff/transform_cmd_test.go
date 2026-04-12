package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTransformEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTransformEnvFile: %v", err)
	}
	return p
}

func TestTransformCmd_MissingArg(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"transform"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestTransformCmd_PrefixAdd(t *testing.T) {
	p := writeTransformEnvFile(t, "HOST=localhost\nPORT=8080\n")
	var buf bytes.Buffer
	transformCmd.SetOut(&buf)
	transformCmd.SetArgs([]string{p, "--prefix-add", "APP_"})
	if err := transformCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_HOST=localhost") {
		t.Errorf("expected APP_HOST in output, got:\n%s", out)
	}
}

func TestTransformCmd_PrefixStrip(t *testing.T) {
	p := writeTransformEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	var buf bytes.Buffer
	transformCmd.SetOut(&buf)
	transformCmd.SetArgs([]string{p, "--prefix-strip", "APP_"})
	if err := transformCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST in output, got:\n%s", out)
	}
}

func TestTransformCmd_UppercaseKeys(t *testing.T) {
	p := writeTransformEnvFile(t, "db_host=localhost\n")
	var buf bytes.Buffer
	transformCmd.SetOut(&buf)
	transformCmd.SetArgs([]string{p, "--uppercase"})
	if err := transformCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got: %s", buf.String())
	}
}

func TestTransformCmd_DropEmpty(t *testing.T) {
	p := writeTransformEnvFile(t, "FOO=bar\nEMPTY=\n")
	var buf bytes.Buffer
	transformCmd.SetOut(&buf)
	transformCmd.SetArgs([]string{p, "--drop-empty"})
	if err := transformCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY to be dropped, got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", out)
	}
}

func TestTransformCmd_InvalidFile(t *testing.T) {
	transformCmd.SetArgs([]string{"/nonexistent/.env"})
	if err := transformCmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}
