package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-conf/internal/envfilter"
	"github.com/yourorg/envoy-conf/internal/envloader"
)

func init() {
	var (
		filterFile    string
		filterPrefix  string
		filterPattern string
		excludeKeys   []string
	)

	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter environment variables by prefix, pattern, or exclusion list",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envloader.LoadFile(filterFile)
			if err != nil {
				return fmt.Errorf("loading env file: %w", err)
			}

			result, err := envfilter.Filter(env, envfilter.Options{
				Prefix:      filterPrefix,
				Pattern:     filterPattern,
				ExcludeKeys: excludeKeys,
			})
			if err != nil {
				return fmt.Errorf("filtering: %w", err)
			}

			keys := make([]string, 0, len(result))
			for k := range result {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result[k])
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d variable(s) matched.\n", len(result))
			return nil
		},
	}

	cmd.Flags().StringVarP(&filterFile, "file", "f", ".env", "Path to the .env file")
	cmd.Flags().StringVar(&filterPrefix, "prefix", "", "Keep only keys with this prefix")
	cmd.Flags().StringVar(&filterPattern, "pattern", "", "Keep only keys matching this regex")
	cmd.Flags().StringArrayVar(&excludeKeys, "exclude", nil, "Exact key names to exclude (repeatable)")

	rootCmd.AddCommand(cmd)
}
