package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-conf/internal/auditor"
)

var (
	auditLogPath string
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "View the audit log for envoy-conf operations",
		RunE:  runAudit,
	}
	auditCmd.Flags().StringVar(&auditLogPath, "log", "audit.json", "path to the audit log file")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	a, err := auditor.LoadLog(auditLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(cmd.OutOrStdout(), "No audit log found.")
			return nil
		}
		return fmt.Errorf("failed to load audit log: %w", err)
	}

	entries := a.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Audit log is empty.")
		return nil
	}

	for _, e := range entries {
		status := "OK"
		if !e.Success {
			status = "FAIL"
		}
		targets := ""
		for i, t := range e.Targets {
			if i > 0 {
				targets += ", "
			}
			targets += t
		}
		line := fmt.Sprintf("[%s] %-10s %-6s %s",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Event,
			status,
			targets,
		)
		if e.Message != "" {
			line += " — " + e.Message
		}
		fmt.Fprintln(cmd.OutOrStdout(), line)
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), a.Summary())
	return nil
}
