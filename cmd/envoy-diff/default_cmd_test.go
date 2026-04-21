package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeDefaultEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestDefaultCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"default"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestDefaultCmd_AppliesMissingKey(t *testing.T) {
	p := writeDefaultEnvFile(t, "A=1\n")
	rootCmd.SetArgs([]string{"default", p, "--set", "B=default_b"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestDefaultCmd_SkipsExistingKey(t *testing.T) {
	p := writeDefaultEnvFile(t, "A=original\n")
	rootCmd.SetArgs([]string{"default", p, "--set", "A=should_not_apply"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestDefaultCmd_OverwriteEmpty(t *testing.T) {
	p := writeDefaultEnvFile(t, "A=\n")
	rootCmd.SetArgs([]string{"default", p, "--set", "A=filled", "--overwrite-empty"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestDefaultCmd_InvalidSetPair(t *testing.T) {
	p := writeDefaultEnvFile(t, "A=1\n")
	rootCmd.SetArgs([]string{"default", p, "--set", "NOEQUALSSIGN"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestSplitKeyValue_Valid(t *testing.T) {
	k, v, ok := splitKeyValue("FOO=bar")
	require.True(t, ok)
	assert.Equal(t, "FOO", k)
	assert.Equal(t, "bar", v)
}

func TestSplitKeyValue_EmptyValue(t *testing.T) {
	k, v, ok := splitKeyValue("FOO=")
	require.True(t, ok)
	assert.Equal(t, "FOO", k)
	assert.Equal(t, "", v)
}

func TestSplitKeyValue_NoEquals(t *testing.T) {
	_, _, ok := splitKeyValue("NOEQUALSSIGN")
	assert.False(t, ok)
}
