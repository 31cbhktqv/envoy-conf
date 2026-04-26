package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-conf/internal/envschema"
)

func writeTempSchema(t *testing.T, schema envschema.Schema) string {
	t.Helper()
	raw, err := json.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "schema-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write(raw)
	_ = f.Close()
	return f.Name()
}

func writeTempSchemaEnv(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func runSchemaCmd(t *testing.T, schemaPath, envPath string) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	rootCmd.SetArgs([]string{"schema", "--schema", schemaPath, "--env", envPath})
	err := rootCmd.Execute()
	_ = w.Close()
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	return string(buf[:n]), err
}

func TestSchemaCmd_AllPass(t *testing.T) {
	schema := envschema.Schema{
		Fields: []envschema.Field{
			{Key: "PORT", Type: envschema.TypeInt, Required: true},
		},
	}
	sp := writeTempSchema(t, schema)
	ep := writeTempSchemaEnv(t, "PORT=9090\n")
	out, err := runSchemaCmd(t, sp, ep)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "passed") {
		t.Errorf("expected passed message, got: %q", out)
	}
}

func TestSchemaCmd_MissingEnvFile(t *testing.T) {
	schema := envschema.Schema{Fields: []envschema.Field{}}
	sp := writeTempSchema(t, schema)
	rootCmd.SetArgs([]string{"schema", "--schema", sp, "--env", "/nonexistent/.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing env file")
	}
}

func TestSchemaCmd_MissingSchemaFile(t *testing.T) {
	ep := writeTempSchemaEnv(t, "PORT=9090\n")
	rootCmd.SetArgs([]string{"schema", "--schema", "/nonexistent/schema.json", "--env", ep})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing schema file")
	}
}
