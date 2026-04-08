package parser_test

import (
	"strings"
	"testing"

	"github.com/envoy-diff/envoy-diff/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEnvFile_BasicKeyValue(t *testing.T) {
	input := `
APP_NAME=envoy-diff
PORT=8080
DEBUG=true
`
	env, err := parser.ParseEnvFile(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, "envoy-diff", env["APP_NAME"])
	assert.Equal(t, "8080", env["PORT"])
	assert.Equal(t, "true", env["DEBUG"])
}

func TestParseEnvFile_IgnoresComments(t *testing.T) {
	input := `
# This is a comment
KEY=value
# Another comment
`
	env, err := parser.ParseEnvFile(strings.NewReader(input))
	require.NoError(t, err)
	assert.Len(t, env, 1)
	assert.Equal(t, "value", env["KEY"])
}

func TestParseEnvFile_QuotedValues(t *testing.T) {
	input := `
DOUBLE="hello world"
SINGLE='foo bar'
`
	env, err := parser.ParseEnvFile(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, "hello world", env["DOUBLE"])
	assert.Equal(t, "foo bar", env["SINGLE"])
}

func TestParseEnvFile_EmptyValue(t *testing.T) {
	input := `EMPTY=`
	env, err := parser.ParseEnvFile(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, "", env["EMPTY"])
}

func TestParseEnvFile_InvalidLine(t *testing.T) {
	input := `NODIVIDER`
	_, err := parser.ParseEnvFile(strings.NewReader(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}

func TestParseEnvFile_EmptyKey(t *testing.T) {
	input := `=value`
	_, err := parser.ParseEnvFile(strings.NewReader(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty key")
}

func TestParseEnvFile_ValueWithEquals(t *testing.T) {
	input := `URL=http://example.com?foo=bar&baz=qux`
	env, err := parser.ParseEnvFile(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, "http://example.com?foo=bar&baz=qux", env["URL"])
}
