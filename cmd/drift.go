package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envdrift"
	"envoy-conf/internal/envloader"
)

var (
	driftTarget    string
	driftBaseFile  string
	driftLiveFile  string
	driftShowMatch bool
	driftColor     bool
)

func init() {
	driftCmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect configuration drift between a baseline and live environment",
		RunE:  runDrift,
	}

	driftCmd.Flags().StringVarP(&driftBaseFile, "baseline", "b", "", "baseline .env file (required)")
	driftCmd.Flags().StringVarP(&driftLiveFile, "live", "l", "", "live .env file (required)")
	driftCmd.Flags().StringVarP(&driftTarget, "target", "t", "unknown", "deployment target name")
	driftCmd.Flags().BoolVar(&driftShowMatch, "show-match", false, "include unchanged keys in output")
	driftCmd.Flags().BoolVar(&driftColor, "color", false, "colorize output")

	_ = driftCmd.MarkFlagRequired("baseline")
	_ = driftCmd.MarkFlagRequired("live")

	rootCmd.AddCommand(driftCmd)
}

func runDrift(cmd *cobra.Command, args []string) error {
	baseline, err := envloader.LoadFile(driftBaseFile)
	if err != nil {
		return fmt.Errorf("loading baseline: %w", err)
	}

	live, err := envloader.LoadFile(driftLiveFile)
	if err != nil {
		return fmt.Errorf("loading live: %w", err)
	}

	report := envdrift.Detect(driftTarget, baseline, live)

	opts := envdrift.RenderOptions{
		Color:     driftColor,
		ShowMatch: driftShowMatch,
	}
	envdrift.Render(os.Stdout, report, opts)

	if report.HasDrift() {
		os.Exit(1)
	}
	return nil
}
