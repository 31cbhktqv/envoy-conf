package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/envrollout"
)

var rolloutStageFiles []string
var rolloutStageNames []string

func init() {
	rolloutCmd := &cobra.Command{
		Use:   "rollout",
		Short: "Plan a multi-stage rollout and check env readiness",
		RunE:  runRollout,
	}
	rolloutCmd.Flags().StringSliceVarP(&rolloutStageFiles, "files", "f", nil, "Ordered env files (e.g. dev.env,staging.env,prod.env)")
	rolloutCmd.Flags().StringSliceVarP(&rolloutStageNames, "names", "n", nil, "Stage names matching --files order")
	_ = rolloutCmd.MarkFlagRequired("files")
	RootCmd.AddCommand(rolloutCmd)
}

func runRollout(cmd *cobra.Command, _ []string) error {
	if len(rolloutStageFiles) < 2 {
		return fmt.Errorf("at least two --files required for a rollout plan")
	}

	stages := make([]envrollout.Stage, 0, len(rolloutStageFiles))
	for i, f := range rolloutStageFiles {
		env, err := envloader.LoadFile(f)
		if err != nil {
			return fmt.Errorf("loading %s: %w", f, err)
		}
		name := f
		if i < len(rolloutStageNames) {
			name = rolloutStageNames[i]
		}
		stages = append(stages, envrollout.Stage{Name: name, Env: env})
	}

	results := envrollout.Plan(stages)
	for _, r := range results {
		status := "✔ ready"
		if !r.Ready {
			status = "✘ BLOCKED"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s → %s  [%s]\n", r.From, r.To, status)
		if len(r.MissingKeys) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  missing : %s\n", strings.Join(r.MissingKeys, ", "))
		}
		if len(r.ChangedKeys) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  changed : %s\n", strings.Join(r.ChangedKeys, ", "))
		}
		if len(r.NewKeys) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  new     : %s\n", strings.Join(r.NewKeys, ", "))
		}
	}

	if envrollout.HasBlocker(results) {
		fmt.Fprintln(cmd.OutOrStdout(), "\nRollout plan has blockers — resolve missing keys before promoting.")
		os.Exit(1)
	}
	return nil
}
