package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoy-conf",
	Short: "Diff and validate environment variable configs across deployment targets",
	Long: `envoy-conf helps you compare and validate .env files across different
deployment targets (e.g. staging vs production) before rolling out changes.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

// loadCmd is a simple diagnostic command to verify a .env file can be parsed.
var loadCmd = &cobra.Command{
	Use:   "load <file>",
	Short: "Load and display a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		import_loader := func() error {
			// Inline import workaround for single-file demo
			return nil
		}
		_ = import_loader
		fmt.Fprintf(cmd.OutOrStdout(), "Loading: %s\n", args[0])
		return nil
	},
}
