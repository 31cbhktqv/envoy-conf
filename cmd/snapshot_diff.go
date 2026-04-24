package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/differ"
	"envoy-conf/internal/envloader"
	"envoy-conf/internal/snapshot"
)

var snapshotDiffNoColor bool

func init() {
	snapshotDiffCmd := &cobra.Command{
		Use:   "snapshot-diff <snapshot-file> <env-file>",
		Short: "Diff a saved snapshot against a current environment file",
		Long:  "Compares a previously saved snapshot to a live env file to surface changes since the snapshot was taken.",
		Args:  cobra.ExactArgs(2),
		RunE:  runSnapshotDiff,
	}

	snapshotDiffCmd.Flags().BoolVar(&snapshotDiffNoColor, "no-color", false, "disable colored output")
	rootCmd.AddCommand(snapshotDiffCmd)
}

func runSnapshotDiff(cmd *cobra.Command, args []string) error {
	snapshotFile := args[0]
	envFile := args[1]

	snap, err := snapshot.Load(snapshotFile)
	if err != nil {
		return fmt.Errorf("failed to load snapshot %q: %w", snapshotFile, err)
	}

	current, err := envloader.LoadFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to load env file %q: %w", envFile, err)
	}

	results := differ.Diff(snap.Env, current)

	opts := differ.DefaultFormatOptions()
	opts.NoColor = snapshotDiffNoColor

	fmt.Fprintf(os.Stdout, "Snapshot target: %s  (taken: %s)\n",
		snap.Target, snap.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(os.Stdout, "Comparing against: %s\n\n", envFile)

	output := differ.Render(results, opts)
	fmt.Fprint(os.Stdout, output)
	fmt.Fprintln(os.Stdout, differ.Summary(results))

	return nil
}
