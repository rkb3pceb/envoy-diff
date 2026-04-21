package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSplitEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeSplitEnvFile: %v", err)
	}
	return p
}

func TestSplitCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"split"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestSplitCmd_InvalidFile(t *testing.T) {
	rootCmd.SetArgs([]string{"split", "/no/such/file.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSplitCmd_SingleValues_PassThrough(t *testing.T) {
	f := writeSplitEnvFile(t, "HOST=localhost\nPORT=8080\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"split", f})
		_ = rootCmd.Execute()
	})
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
}

func TestSplitCmd_SplitsCommaList(t *testing.T) {
	f := writeSplitEnvFile(t, "HOSTS=a,b,c\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"split", f})
		_ = rootCmd.Execute()
	})
	if !strings.Contains(out, "HOSTS_1=a") {
		t.Errorf("expected HOSTS_1=a in output, got: %s", out)
	}
	if !strings.Contains(out, "HOSTS_2=b") {
		t.Errorf("expected HOSTS_2=b in output, got: %s", out)
	}
	if !strings.Contains(out, "HOSTS_3=c") {
		t.Errorf("expected HOSTS_3=c in output, got: %s", out)
	}
}

func TestSplitCmd_CustomDelimiter(t *testing.T) {
	f := writeSplitEnvFile(t, "PORTS=80|443\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"split", "--delimiter", "|", f})
		_ = rootCmd.Execute()
	})
	if !strings.Contains(out, "PORTS_1=80") {
		t.Errorf("expected PORTS_1=80, got: %s", out)
	}
	if !strings.Contains(out, "PORTS_2=443") {
		t.Errorf("expected PORTS_2=443, got: %s", out)
	}
}

func TestSplitCmd_KeepOriginalFlag(t *testing.T) {
	f := writeSplitEnvFile(t, "TAGS=x,y\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"split", "--keep-original", f})
		_ = rootCmd.Execute()
	})
	if !strings.Contains(out, "TAGS=x,y") {
		t.Errorf("expected original TAGS key preserved, got: %s", out)
	}
	if !strings.Contains(out, "TAGS_1=x") {
		t.Errorf("expected TAGS_1=x, got: %s", out)
	}
}
