package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/differ"
	"envoy-conf/internal/envloader"
)

var (
	diffLabelA string
	diffLabelB string
	diffNoColor bool
	diffVerbose bool
)

var diffCmd = &cobra.Command{
	Use:   "diff <file-a> <file-b>",
	Short: "Diff two .env files and show differences",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileA, fileB := args[0], args[1]

		envA, err := envloader.LoadFile(fileA)
		if err != nil {
			return fmt.Errorf("loading %s: %w", fileA, err)
		}

		envB, err := envloader.LoadFile(fileB)
		if err != nil {
			return fmt.Errorf("loading %s: %w", fileB, err)
		}

		result := differ.Diff(envA, envB)

		opts := differ.FormatOptions{
			Color:   !diffNoColor,
			Verbose: diffVerbose,
			LabelA:  diffLabelA,
			LabelB:  diffLabelB,
		}
		if opts.LabelA == "" {
			opts.LabelA = fileA
		}
		if opts.LabelB == "" {
			opts.LabelB = fileB
		}

		differ.Render(os.Stdout, result, opts)
		fmt.Fprintf(os.Stderr, "Summary: %s\n", differ.Summary(result))

		if result.HasDifferences() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().StringVar(&diffLabelA, "label-a", "", "label for file-a (default: filename)")
	diffCmd.Flags().StringVar(&diffLabelB, "label-b", "", "label for file-b (default: filename)")
	diffCmd.Flags().BoolVar(&diffNoColor, "no-color", false, "disable colored output")
	diffCmd.Flags().BoolVarP(&diffVerbose, "verbose", "v", false, "show unchanged keys as well")
	rootCmd.AddCommand(diffCmd)
}
