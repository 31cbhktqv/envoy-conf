package envloader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	env, err := LoadFile(path)
	require.NoError(t, err)
	assert.Equal(t, EnvMap{"FOO": "bar", "BAZ": "qux"}, env)
}

func TestLoadFile_CommentsAndBlanks(t *testing.T) {
	content := "# this is a comment\n\nKEY=value\n"
	path := writeTempEnv(t, content)
	env, err := LoadFile(path)
	require.NoError(t, err)
	assert.Equal(t, EnvMap{"KEY": "value"}, env)
}

func TestLoadFile_InvalidSyntax(t *testing.T) {
	path := writeTempEnv(t, "NOEQUALSIGN\n")
	_, err := LoadFile(path)
	assert.ErrorContains(t, err, "invalid syntax")
}

func TestLoadFile_EmptyKey(t *testing.T) {
	path := writeTempEnv(t, "=value\n")
	_, err := LoadFile(path)
	assert.ErrorContains(t, err, "empty key")
}

func TestLoadFile_FileNotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path.env")
	assert.ErrorContains(t, err, "opening env file")
}

func TestLoadFile_ValueWithEquals(t *testing.T) {
	path := writeTempEnv(t, "URL=http://example.com?foo=bar\n")
	env, err := LoadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "http://example.com?foo=bar", env["URL"])
}
