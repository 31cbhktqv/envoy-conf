package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/envschema"
)

var (
	schemaFile  string
	schemaEnv   string
	schemaQuiet bool
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate an env file against a JSON schema definition",
		RunE:  runSchema,
	}

	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "path to JSON schema file (required)")
	schemaCmd.Flags().StringVarP(&schemaEnv, "env", "e", "", "path to .env file (required)")
	schemaCmd.Flags().BoolVarP(&schemaQuiet, "quiet", "q", false, "suppress output, use exit code only")
	_ = schemaCmd.MarkFlagRequired("schema")
	_ = schemaCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	raw, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var schema envschema.Schema
	if err := json.Unmarshal(raw, &schema); err != nil {
		return fmt.Errorf("parsing schema JSON: %w", err)
	}

	env, err := envloader.LoadFile(schemaEnv)
	if err != nil {
		return fmt.Errorf("loading env file: %w", err)
	}

	violations := schema.Validate(env)

	if schemaQuiet {
		if len(violations) > 0 {
			os.Exit(1)
		}
		return nil
	}

	if len(violations) == 0 {
		fmt.Println("✔ schema validation passed — no violations found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "✖ schema validation failed (%d violation(s)):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", v.Key, v.Message)
	}
	os.Exit(1)
	return nil
}
