package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envcanary"
	"envoy-conf/internal/envloader"
)

var (
	canaryBaselineFile  string
	canaryCurrentFile   string
	canaryRequiredKeys  []string
	canaryWatchKeys     []string
	canaryAllowMissing  bool
	canaryShowOK        bool
)

func init() {
	canaryCmd := &cobra.Command{
		Use:   "canary",
		Short: "Check environment variables against a baseline for safe rollout",
		RunE:  runCanary,
	}

	canaryCmd.Flags().StringVarP(&canaryBaselineFile, "baseline", "b", "", "baseline .env file (required)")
	canaryCmd.Flags().StringVarP(&canaryCurrentFile, "current", "c", "", "current .env file (required)")
	canaryCmd.Flags().StringSliceVar(&canaryRequiredKeys, "require", nil, "keys that must be present in current env")
	canaryCmd.Flags().StringSliceVar(&canaryWatchKeys, "watch", nil, "keys whose value changes should be flagged")
	canaryCmd.Flags().BoolVar(&canaryAllowMissing, "allow-missing", false, "downgrade missing required keys from critical to warning")
	canaryCmd.Flags().BoolVar(&canaryShowOK, "show-ok", false, "include passing checks in output")

	_ = canaryCmd.MarkFlagRequired("baseline")
	_ = canaryCmd.MarkFlagRequired("current")

	rootCmd.AddCommand(canaryCmd)
}

func runCanary(cmd *cobra.Command, _ []string) error {
	baseline, err := envloader.LoadFile(canaryBaselineFile)
	if err != nil {
		return fmt.Errorf("loading baseline: %w", err)
	}
	current, err := envloader.LoadFile(canaryCurrentFile)
	if err != nil {
		return fmt.Errorf("loading current: %w", err)
	}

	// Support KEY=VALUE pairs passed directly via --require / --watch flags
	requiredKeys := normaliseKeyList(canaryRequiredKeys)
	watchKeys := normaliseKeyList(canaryWatchKeys)

	opts := envcanary.Options{
		RequiredKeys: requiredKeys,
		WatchKeys:    watchKeys,
		AllowMissing: canaryAllowMissing,
	}

	results := envcanary.Check(baseline, current, opts)
	envcanary.Render(cmd.OutOrStdout(), results, canaryShowOK)

	if envcanary.HasCritical(results) {
		os.Exit(1)
	}
	return nil
}

// normaliseKeyList strips any accidental KEY=VALUE formatting from flag values.
func normaliseKeyList(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if idx := strings.Index(k, "="); idx != -1 {
			k = k[:idx]
		}
		k = strings.TrimSpace(k)
		if k != "" {
			out = append(out, k)
		}
	}
	return out
}
