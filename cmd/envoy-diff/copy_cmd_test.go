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

func writeCopyEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestCopyCmd_MissingArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"copy"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestCopyCmd_CopiesAllKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnvFile(t, dir, "src.env", "FOO=1\nBAR=2\n")
	dst := writeCopyEnvFile(t, dir, "dst.env", "")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"copy", src, dst})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "copied 2")

	data, _ := os.ReadFile(dst)
	assert.Contains(t, string(data), "FOO=1")
	assert.Contains(t, string(data), "BAR=2")
}

func TestCopyCmd_DestPrefixApplied(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnvFile(t, dir, "src.env", "KEY=val\n")
	dst := writeCopyEnvFile(t, dir, "dst.env", "")

	rootCmd.SetArgs([]string{"copy", "--dest-prefix", "PROD_", src, dst})
	err := rootCmd.Execute()
	require.NoError(t, err)

	data, _ := os.ReadFile(dst)
	assert.Contains(t, string(data), "PROD_KEY=val")
}

func TestCopyCmd_SkipsExistingWithoutOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnvFile(t, dir, "src.env", "X=new\n")
	dst := writeCopyEnvFile(t, dir, "dst.env", "X=old\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"copy", src, dst})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "skipped 1")

	data, _ := os.ReadFile(dst)
	assert.Contains(t, string(data), "X=old")
}

func TestCopyCmd_OverwriteFlag(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnvFile(t, dir, "src.env", "X=new\n")
	dst := writeCopyEnvFile(t, dir, "dst.env", "X=old\n")

	rootCmd.SetArgs([]string{"copy", "--overwrite", src, dst})
	err := rootCmd.Execute()
	require.NoError(t, err)

	data, _ := os.ReadFile(dst)
	assert.True(t, strings.Contains(string(data), "X=new"))
}

func TestCopyCmd_MissingDstFile_CreatesNew(t *testing.T) {
	dir := t.TempDir()
	src := writeCopyEnvFile(t, dir, "src.env", "NEW=1\n")
	dst := filepath.Join(dir, "nonexistent.env")

	rootCmd.SetArgs([]string{"copy", src, dst})
	err := rootCmd.Execute()
	require.NoError(t, err)

	data, _ := os.ReadFile(dst)
	assert.Contains(t, string(data), "NEW=1")
}
