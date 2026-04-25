package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/resolver"
)

var (
	resolveFiles      []string
	resolveOSFallback bool
	resolveOverrides  []string
)

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Merge multiple env files into a single resolved config",
	Long:  "Loads one or more .env files in order and merges them into a single flat config, with optional OS environment fallback and key overrides.",
	RunE:  runResolve,
}

func init() {
	rootCmd.AddCommand(resolveCmd)
	resolveCmd.Flags().StringArrayVarP(&resolveFiles, "file", "f", nil, "Env file(s) to resolve (applied in order)")
	resolveCmd.Flags().BoolVar(&resolveOSFallback, "os-fallback", false, "Use OS environment as base layer")
	resolveCmd.Flags().StringArrayVar(&resolveOverrides, "override-key", nil, "Keys to override from OS environment after merge")
	_ = resolveCmd.MarkFlagRequired("file")
}

func runResolve(cmd *cobra.Command, args []string) error {
	sources, err := loadSources(resolveFiles)
	if err != nil {
		return err
	}

	opts := resolver.ResolveOptions{
		FallbackToOS: resolveOSFallback,
		OverrideKeys: resolveOverrides,
	}

	result, err := resolver.Resolve(sources, opts)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(cmd.OutOrStdout(), "# Resolved from: %v\n", resolver.SourceNames(sources))
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result[k])
	}
	return nil
}

// loadSources reads each file path into a resolver.Source, returning the first
// error encountered along with the path that caused it.
func loadSources(paths []string) ([]resolver.Source, error) {
	sources := make([]resolver.Source, 0, len(paths))
	for _, path := range paths {
		vars, err := envloader.LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load %q: %w", path, err)
		}
		sources = append(sources, resolver.Source{Name: path, Vars: vars})
	}
	return sources, nil
}
