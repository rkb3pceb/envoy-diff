package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGroupEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGroupCmd_MissingArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "group")
	if err == nil {
		t.Fatal("expected error for missing argument")
	}
}

func TestGroupCmd_InvalidFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "group", "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestGroupCmd_SingleGroup(t *testing.T) {
	p := writeGroupEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	out, err := executeCommand(rootCmd, "group", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] group in output, got:\n%s", out)
	}
}

func TestGroupCmd_MultipleGroups(t *testing.T) {
	p := writeGroupEnvFile(t, "DB_HOST=localhost\nAPP_PORT=8080\nREDIS_URL=redis://localhost\n")
	out, err := executeCommand(rootCmd, "group", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, label := range []string{"[APP]", "[DB]", "[REDIS]"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected %s in output, got:\n%s", label, out)
		}
	}
}

func TestGroupCmd_NoOtherFlag_ExcludesUngrouped(t *testing.T) {
	p := writeGroupEnvFile(t, "HOSTNAME=box1\nDB_HOST=localhost\n")
	out, err := executeCommand(rootCmd, "group", "--no-other", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "OTHER") {
		t.Errorf("did not expect OTHER group, got:\n%s", out)
	}
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] group, got:\n%s", out)
	}
}

func TestGroupCmd_CustomDelimiter(t *testing.T) {
	p := writeGroupEnvFile(t, "DB.HOST=localhost\nDB.PORT=5432\n")
	out, err := executeCommand(rootCmd, "group", "--delimiter=.", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] group with dot delimiter, got:\n%s", out)
	}
}
