package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/linter"
)

var lintFailOnWarn bool

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint [file]",
		Short: "Lint an env file for common issues",
		Args:  cobra.ExactArgs(1),
		RunE:  runLint,
	}
	lintCmd.Flags().BoolVar(&lintFailOnWarn, "fail", false, "exit with non-zero status if violations are found")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	path := args[0]

	env, err := envloader.LoadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load file %q: %w", path, err)
	}

	l := linter.New(linter.DefaultRules())
	violations := l.Lint(env)

	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ No linting violations found.")
		return nil
	}

	// Sort violations for deterministic output
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].Key != violations[j].Key {
			return violations[i].Key < violations[j].Key
		}
		return violations[i].Rule < violations[j].Rule
	})

	fmt.Fprintf(cmd.OutOrStdout(), "⚠ %d linting violation(s) in %s:\n\n", len(violations), path)
	for _, v := range violations {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s\n", v.Rule, v.Message)
	}

	if lintFailOnWarn {
		os.Exit(1)
	}
	return nil
}
