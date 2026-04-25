package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempTemplateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempTemplateEnv: %v", err)
	}
	return p
}

func runTemplateCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"template"}, args...))
	err := rootCmd.Execute()
	// Reset for subsequent tests
	rootCmd.SetArgs([]string{})
	return buf.String(), err
}

func TestTemplateCmd_BasicExpansion(t *testing.T) {
	envFile := writeTempTemplateEnv(t, "DSN=postgres://${HOST}:5432/db\n")
	lookupFile := writeTempTemplateEnv(t, "HOST=localhost\n")

	out, err := runTemplateCmd(t, []string{"--file", envFile, "--lookup", lookupFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DSN=postgres://localhost:5432/db") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTemplateCmd_NoPlaceholders(t *testing.T) {
	envFile := writeTempTemplateEnv(t, "FOO=bar\nBAZ=qux\n")
	out, err := runTemplateCmd(t, []string{"--file", envFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTemplateCmd_StrictMissingVar(t *testing.T) {
	envFile := writeTempTemplateEnv(t, "KEY=${UNDEFINED_XYZ}\n")
	_, err := runTemplateCmd(t, []string{"--file", envFile, "--strict"})
	if err == nil {
		t.Fatal("expected error in strict mode for unresolved variable")
	}
}

func TestTemplateCmd_MissingFile(t *testing.T) {
	_ = &cobra.Command{} // ensure import used
	_, err := runTemplateCmd(t, []string{"--file", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
