package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCastEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return p
}

func TestCastCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"cast"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestCastCmd_StringType_AllKeys(t *testing.T) {
	dir := t.TempDir()
	f := writeCastEnvFile(t, dir, ".env", "HOST=localhost\nPORT=8080\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"cast", "--type", "string", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("expected ok status in output, got: %s", out)
	}
}

func TestCastCmd_IntType_ValidValues(t *testing.T) {
	dir := t.TempDir()
	f := writeCastEnvFile(t, dir, ".env", "PORT=9090\nTIMEOUT=30\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"cast", "--type", "int", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ok") {
		t.Errorf("expected ok status, got: %s", out)
	}
}

func TestCastCmd_BoolType_Normalises(t *testing.T) {
	dir := t.TempDir()
	f := writeCastEnvFile(t, dir, ".env", "ENABLED=yes\nVERBOSE=0\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"cast", "--type", "bool", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ENABLED") {
		t.Errorf("expected ENABLED in output, got: %s", out)
	}
}

func TestCastCmd_IntType_InvalidValue_SkipErrors(t *testing.T) {
	dir := t.TempDir()
	f := writeCastEnvFile(t, dir, ".env", "PORT=not-a-number\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"cast", "--type", "int", "--skip-errors", f})
	// should not return error when skip-errors is set
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error with skip-errors: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL status in output, got: %s", out)
	}
}

func TestCastCmd_SpecificKey_OnlyCastsThat(t *testing.T) {
	dir := t.TempDir()
	f := writeCastEnvFile(t, dir, ".env", "PORT=8080\nNAME=hello\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"cast", "--type", "int", "--key", "PORT", f})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "NAME") {
		t.Errorf("NAME should not appear when specific key requested, got: %s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("PORT should appear in output, got: %s", out)
	}
}
