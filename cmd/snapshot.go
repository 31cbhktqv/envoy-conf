package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/snapshot"
)

var (
	snapshotTarget string
	snapshotOutput string
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot <env-file>",
		Short: "Save a snapshot of an environment file for later comparison",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshot,
	}

	snapshotCmd.Flags().StringVarP(&snapshotTarget, "target", "t", "", "deployment target label (required)")
	snapshotCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "", "output file path for the snapshot (required)")
	snapshotCmd.MarkFlagRequired("target")
	snapshotCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	envFile := args[0]

	env, err := envloader.LoadFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to load env file %q: %w", envFile, err)
	}

	s := snapshot.New(snapshotTarget, env)

	if err := snapshot.Save(s, snapshotOutput); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot saved: target=%s keys=%d file=%s\n",
		s.Target, len(s.Env), snapshotOutput)
	return nil
}
