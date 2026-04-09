package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootCmd_RequiresTwoArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"file1.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when only one argument provided")
	}
}

func TestRunDiff_InvalidOldFile(t *testing.T) {
	tmpDir := t.TempDir()
	newFile := filepath.Join(tmpDir, "new.env")
	os.WriteFile(newFile, []byte("KEY=value\n"), 0644)

	err := runDiff(rootCmd, []string{"nonexistent.env", newFile})
	if err == nil {
		t.Error("expected error for nonexistent old file")
	}
}

func TestRunDiff_InvalidNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.env")
	os.WriteFile(oldFile, []byte("KEY=value\n"), 0644)

	err := runDiff(rootCmd, []string{oldFile, "nonexistent.env"})
	if err == nil {
		t.Error("expected error for nonexistent new file")
	}
}

func TestRunDiff_Success(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.env")
	newFile := filepath.Join(tmpDir, "new.env")

	os.WriteFile(oldFile, []byte("KEY=old\n"), 0644)
	os.WriteFile(newFile, []byte("KEY=new\n"), 0644)

	// Reset flags
	format = "text"
	auditMode = false
	noColor = true

	err := runDiff(rootCmd, []string{oldFile, newFile})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunDiff_WithAuditMode(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.env")
	newFile := filepath.Join(tmpDir, "new.env")

	os.WriteFile(oldFile, []byte("API_KEY=old_secret\n"), 0644)
	os.WriteFile(newFile, []byte("API_KEY=new_secret\n"), 0644)

	// Reset flags
	format = "json"
	auditMode = true
	noColor = true

	err := runDiff(rootCmd, []string{oldFile, newFile})
	// Should return error because audit found issues
	if err == nil {
		t.Error("expected error from audit findings")
	}
}

func TestRunDiff_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.env")
	newFile := filepath.Join(tmpDir, "new.env")

	os.WriteFile(oldFile, []byte("KEY=value\n"), 0644)
	os.WriteFile(newFile, []byte("KEY=value\n"), 0644)

	format = "invalid"
	auditMode = false

	err := runDiff(rootCmd, []string{oldFile, newFile})
	if err == nil {
		t.Error("expected error for invalid format")
	}
}
