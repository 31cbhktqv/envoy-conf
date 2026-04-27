package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/envpromote"
)

var promoteTo string

func init() {
	promoteCmd := &cobra.Command{
		Use:   "promote <from-file> <to-file>",
		Short: "Check whether an environment is ready to be promoted to another stage",
		Args:  cobra.ExactArgs(2),
		RunE:  runPromote,
	}

	promoteCmd.Flags().StringVar(&promoteTo, "from-name", "source", "label for the source stage")
	promoteCmd.Flags().StringVar(&promoteTo, "to-name", "target", "label for the target stage")

	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	fromName, _ := cmd.Flags().GetString("from-name")
	toName, _ := cmd.Flags().GetString("to-name")

	fromEnv, err := envloader.LoadFile(args[0])
	if err != nil {
		return fmt.Errorf("loading source file: %w", err)
	}

	toEnv, err := envloader.LoadFile(args[1])
	if err != nil {
		return fmt.Errorf("loading target file: %w", err)
	}

	from := envpromote.Stage{Name: fromName, Env: fromEnv}
	to := envpromote.Stage{Name: toName, Env: toEnv}

	result := envpromote.Promote(from, to)
	envpromote.Render(os.Stdout, result)

	if !result.Ready {
		os.Exit(1)
	}
	return nil
}
