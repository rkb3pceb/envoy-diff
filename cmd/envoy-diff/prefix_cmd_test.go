package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePrefixEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return p
}

func TestPrefixCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"prefix"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestPrefixCmd_NoFlagError(t *testing.T) {
	f := writePrefixEnvFile(t, "FOO=1\n")
	rootCmd.SetArgs([]string{"prefix", f})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error when neither --add nor --strip given")
	}
}

func TestPrefixCmd_AddPrefix(t *testing.T) {
	f := writePrefixEnvFile(t, "FOO=1\nBAR=2\n")
	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"prefix", "--add", "APP_", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_FOO=1") {
		t.Errorf("expected APP_FOO=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "APP_BAR=2") {
		t.Errorf("expected APP_BAR=2 in output, got: %s", out)
	}
}

func TestPrefixCmd_StripPrefix(t *testing.T) {
	f := writePrefixEnvFile(t, "APP_FOO=1\nAPP_BAR=2\nOTHER=3\n")
	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"prefix", "--strip", "APP_", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO=1") {
		t.Errorf("expected FOO=1 in output, got: %s", out)
	}
	if strings.Contains(out, "APP_FOO") {
		t.Errorf("APP_FOO should have been stripped, got: %s", out)
	}
	if !strings.Contains(out, "OTHER=3") {
		t.Errorf("expected OTHER=3 to pass through, got: %s", out)
	}
}

func TestPrefixCmd_NoDouble_PreventsDuplication(t *testing.T) {
	f := writePrefixEnvFile(t, "APP_FOO=1\nBAR=2\n")
	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"prefix", "--add", "APP_", "--no-double", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "APP_APP_FOO") {
		t.Errorf("should not double-prefix, got: %s", out)
	}
	if !strings.Contains(out, "APP_FOO=1") {
		t.Errorf("expected APP_FOO=1, got: %s", out)
	}
}
