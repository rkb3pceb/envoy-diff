package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeRotateEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestRotateCmd_MissingArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "rotate")
	assert.Error(t, err)
}

func TestRotateCmd_NoRuleFlagError(t *testing.T) {
	p := writeRotateEnvFile(t, "A=1\n")
	_, err := executeCommand(rootCmd, "rotate", p)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--rule")
}

func TestRotateCmd_InvalidRuleFormat(t *testing.T) {
	p := writeRotateEnvFile(t, "A=1\n")
	_, err := executeCommand(rootCmd, "rotate", "--rule", "BADFORMAT", p)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid rule")
}

func TestRotateCmd_RenamesKey(t *testing.T) {
	p := writeRotateEnvFile(t, "OLD_KEY=hello\nOTHER=world\n")
	out, err := executeCommand(rootCmd, "rotate", "--rule", "OLD_KEY=NEW_KEY", p)
	require.NoError(t, err)
	assert.Contains(t, out, "NEW_KEY=hello")
	assert.NotContains(t, out, "OLD_KEY=")
	assert.Contains(t, out, "OTHER=world")
}

func TestRotateCmd_KeepOldRetainsBothKeys(t *testing.T) {
	p := writeRotateEnvFile(t, "OLD_KEY=hello\n")
	out, err := executeCommand(rootCmd, "rotate", "--rule", "OLD_KEY=NEW_KEY", "--keep-old", p)
	require.NoError(t, err)
	assert.Contains(t, out, "NEW_KEY=hello")
	assert.Contains(t, out, "OLD_KEY=hello")
}

func TestRotateCmd_ErrorOnMissing(t *testing.T) {
	p := writeRotateEnvFile(t, "A=1\n")
	_, err := executeCommand(rootCmd, "rotate", "--rule", "MISSING=NEW", "--error-on-missing", p)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MISSING")
}

// executeCommand is assumed to be defined in main_test.go as a shared test helper.
var _ = bytes.NewBuffer // ensure bytes is used
var _ = strings.Contains
