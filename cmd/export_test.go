package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempExportEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func runExportCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var sb strings.Builder
	rootCmd.SetOut(&sb)
	rootCmd.SetErr(&sb)
	rootCmd.SetArgs(append([]string{"export"}, args...))
	_, err := rootCmd.ExecuteC()
	return sb.String(), err
}

func resetExportFlags() {
	exportFormat = "dotenv"
	exportMaskKeys = nil
	exportAutoMask = false
	exportOutputFile = ""
}

func TestExportCmd_DotenvDefault(t *testing.T) {
	defer resetExportFlags()
	p := writeTempExportEnv(t, "APP_ENV=staging\nPORT=8080\n")
	out, err := runExportCmd(t, p)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "APP_ENV=staging") {
		t.Errorf("expected APP_ENV in output, got: %s", out)
	}
}

func TestExportCmd_JSONFormat(t *testing.T) {
	defer resetExportFlags()
	p := writeTempExportEnv(t, "KEY=value\n")
	out, err := runExportCmd(t, "--format", "json", p)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestExportCmd_MaskFlag(t *testing.T) {
	defer resetExportFlags()
	p := writeTempExportEnv(t, "DB_PASSWORD=supersecret\nAPP=myapp\n")
	out, err := runExportCmd(t, "--mask", "DB_PASSWORD", p)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
}

func TestExportCmd_OutputFile(t *testing.T) {
	defer resetExportFlags()
	p := writeTempExportEnv(t, "FOO=bar\n")
	outFile := filepath.Join(t.TempDir(), "out.env")
	_, err := runExportCmd(t, "--output", outFile, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("expected FOO=bar in output file, got: %s", data)
	}
}

func TestExportCmd_MissingFile(t *testing.T) {
	defer resetExportFlags()
	// Suppress cobra usage output
	rootCmd.SetOut(os.Discard)
	rootCmd.SetErr(os.Discard)
	rootCmd.SetArgs([]string{"export", "/nonexistent/.env"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Error("expected error for missing file")
	}
}

var _ = func() *cobra.Command { return rootCmd }()
