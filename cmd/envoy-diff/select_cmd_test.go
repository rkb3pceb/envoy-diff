package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeSelectEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestSelectCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"select"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestSelectCmd_InvalidFile(t *testing.T) {
	rootCmd.SetArgs([]string{"select", "/no/such/file.env"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestSelectCmd_SelectSpecificKey(t *testing.T) {
	f := writeSelectEnvFile(t, "FOO=bar\nBAR=baz\nBAZ=qux\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"select", "--key", "FOO", f})
		require.NoError(t, rootCmd.Execute())
	})
	assert.Contains(t, out, "FOO=bar")
	assert.NotContains(t, out, "BAR")
	assert.NotContains(t, out, "BAZ")
}

func TestSelectCmd_SelectByPrefix(t *testing.T) {
	f := writeSelectEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=pg://\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"select", "--prefix", "APP_", f})
		require.NoError(t, rootCmd.Execute())
	})
	assert.Contains(t, out, "APP_HOST=localhost")
	assert.Contains(t, out, "APP_PORT=8080")
	assert.NotContains(t, out, "DB_URL")
}

func TestSelectCmd_Invert(t *testing.T) {
	f := writeSelectEnvFile(t, "FOO=1\nBAR=2\nBAZ=3\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"select", "--key", "FOO", "--invert", f})
		require.NoError(t, rootCmd.Execute())
	})
	assert.NotContains(t, out, "FOO")
	assert.True(t, strings.Contains(out, "BAR") || strings.Contains(out, "BAZ"))
}

func TestSelectCmd_CaseInsensitive(t *testing.T) {
	f := writeSelectEnvFile(t, "FOO=hello\nBAR=world\n")
	out := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"select", "--key", "foo", "--case-insensitive", f})
		require.NoError(t, rootCmd.Execute())
	})
	assert.Contains(t, out, "FOO=hello")
	assert.NotContains(t, out, "BAR")
}
