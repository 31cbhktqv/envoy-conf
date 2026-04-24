package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func runLintCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	// Re-use runLint directly via a fresh cobra command
	cmd := &cobra.Command{Use: "lint", RunE: runLint}
	cmd.Flags().BoolVar(&lintFailOnWarn, "fail", false, "")
	cmd.SetOut(buf)
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	return buf.String(), err
}

func TestLintCmd_NoViolations(t *testing.T) {
	path := writeTempLintEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	out, err := runLintCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No linting violations") {
		t.Errorf("expected clean message, got: %s", out)
	}
}

func TestLintCmd_WithViolations(t *testing.T) {
	path := writeTempLintEnv(t, "app_host=localhost\nAPP_SECRET=\n")
	out, err := runLintCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "violation") {
		t.Errorf("expected violations in output, got: %s", out)
	}
	if !strings.Contains(out, "no-lowercase-key") {
		t.Errorf("expected no-lowercase-key rule in output, got: %s", out)
	}
}

func TestLintCmd_MissingFile(t *testing.T) {
	_, err := runLintCmd(t, "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
