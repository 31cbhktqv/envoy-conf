package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envgraph"
	"envoy-conf/internal/envloader"
)

var (
	graphDOT      bool
	graphHideDeps bool
)

func init() {
	graphCmd := &cobra.Command{
		Use:   "graph <env-file>",
		Short: "Visualise dependency relationships between environment variables",
		Args:  cobra.ExactArgs(1),
		RunE:  runGraph,
	}

	graphCmd.Flags().BoolVar(&graphDOT, "dot", false, "Output Graphviz DOT format")
	graphCmd.Flags().BoolVar(&graphHideDeps, "no-deps", false, "Hide dependency annotations in plain output")

	rootCmd.AddCommand(graphCmd)
}

func runGraph(cmd *cobra.Command, args []string) error {
	env, err := envloader.LoadFile(args[0])
	if err != nil {
		return fmt.Errorf("loading %q: %w", args[0], err)
	}

	g := envgraph.Build(env)

	if graphDOT {
		return envgraph.RenderDOT(os.Stdout, g)
	}

	opts := envgraph.DefaultRenderOptions()
	if graphHideDeps {
		opts.ShowDeps = false
	}
	return envgraph.Render(os.Stdout, g, opts)
}
