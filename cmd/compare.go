package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-conf/internal/envcompare"
	"github.com/yourorg/envoy-conf/internal/envloader"
)

func init() {
	var ignoreKeys []string
	var ignorePatterns []string
	var caseInsensitive bool

	cmd := &cobra.Command{
		Use:   "compare <fileA> <fileB>",
		Short: "Compare two env files and report differences",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := envloader.LoadFile(args[0])
			if err != nil {
				return fmt.Errorf("loading %s: %w", args[0], err)
			}
			b, err := envloader.LoadFile(args[1])
			if err != nil {
				return fmt.Errorf("loading %s: %w", args[1], err)
			}

			opts := envcompare.CompareOptions{
				IgnoreKeys:          ignoreKeys,
				IgnorePatterns:      ignorePatterns,
				CaseSensitiveValues: !caseInsensitive,
			}

			res, err := envcompare.Compare(a, b, opts)
			if err != nil {
				return err
			}

			hasIssues := len(res.MissingInA)+len(res.MissingInB)+len(res.Mismatched) > 0

			if !hasIssues {
				fmt.Fprintf(cmd.OutOrStdout(), "✔ No differences found (%d keys matched)\n", res.MatchedCount)
				return nil
			}

			sort.Strings(res.MissingInB)
			for _, k := range res.MissingInB {
				fmt.Fprintf(cmd.OutOrStdout(), "- %-30s  (only in %s)\n", k, args[0])
			}

			sort.Strings(res.MissingInA)
			for _, k := range res.MissingInA {
				fmt.Fprintf(cmd.OutOrStdout(), "+ %-30s  (only in %s)\n", k, args[1])
			}

			mismatchKeys := make([]string, 0, len(res.Mismatched))
			for k := range res.Mismatched {
				mismatchKeys = append(mismatchKeys, k)
			}
			sort.Strings(mismatchKeys)
			for _, k := range mismatchKeys {
				pair := res.Mismatched[k]
				fmt.Fprintf(cmd.OutOrStdout(), "~ %-30s  %q → %q\n", k, pair[0], pair[1])
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nSummary: %d missing in B, %d missing in A, %d mismatched, %d matched\n",
				len(res.MissingInB), len(res.MissingInA), len(res.Mismatched), res.MatchedCount)

			os.Exit(1)
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&ignoreKeys, "ignore-key", nil, "exact key names to ignore")
	cmd.Flags().StringArrayVar(&ignorePatterns, "ignore-pattern", nil, "regex patterns for keys to ignore")
	cmd.Flags().BoolVar(&caseInsensitive, "case-insensitive", false, "treat values as case-insensitive")

	rootCmd.AddCommand(cmd)
}
