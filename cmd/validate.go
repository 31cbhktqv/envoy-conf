package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/validator"
)

var rulesFile string

var validateCmd = &cobra.Command{
	Use:   "validate [env-file]",
	Short: "Validate an env file against a set of rules",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile := args[0]

		env, err := envloader.LoadFile(envFile)
		if err != nil {
			return fmt.Errorf("loading env file: %w", err)
		}

		var rules []validator.Rule
		if rulesFile != "" {
			data, err := os.ReadFile(rulesFile)
			if err != nil {
				return fmt.Errorf("reading rules file: %w", err)
			}
			if err := json.Unmarshal(data, &rules); err != nil {
				return fmt.Errorf("parsing rules file: %w", err)
			}
		}

		violations := validator.Validate(env, rules)
		if len(violations) == 0 {
			fmt.Println("✓ All validation rules passed.")
			return nil
		}

		fmt.Fprintf(os.Stderr, "✗ %d violation(s) found:\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "  - %s\n", v.Error())
		}
		os.Exit(1)
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "Path to JSON rules file")
	rootCmd.AddCommand(validateCmd)
}
