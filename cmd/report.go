package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/envoy-conf/internal/differ"
	"github.com/user/envoy-conf/internal/envloader"
	"github.com/user/envoy-conf/internal/reporter"
	"github.com/user/envoy-conf/internal/validator"
)

var (
	reportFormat string
	reportOutput string
	reportRules  string
)

func init() {
	reportCmd := &cobra.Command{
		Use:   "report <source.env> <target.env>",
		Short: "Generate a combined diff and validation report",
		Args:  cobra.ExactArgs(2),
		RunE:  runReport,
	}

	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "text", "Output format: text or json")
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "", "Write report to file instead of stdout")
	reportCmd.Flags().StringVarP(&reportRules, "rules", "r", "", "Path to validation rules file (optional)")

	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	sourceFile, targetFile := args[0], args[1]

	source, err := envloader.LoadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("loading source: %w", err)
	}
	target, err := envloader.LoadFile(targetFile)
	if err != nil {
		return fmt.Errorf("loading target: %w", err)
	}

	diffItems := differ.Diff(source, target)
	formatOpts := differ.DefaultFormatOptions()
	summary := differ.Summary(diffItems)

	var violations []validator.Violation
	if reportRules != "" {
		rules, rerr := validator.LoadRules(reportRules)
		if rerr != nil {
			return fmt.Errorf("loading rules: %w", rerr)
		}
		violations = validator.Validate(target, rules)
	}

	r := reporter.Report{
		GeneratedAt: time.Now().UTC(),
		SourceFile:  sourceFile,
		TargetFile:  targetFile,
		DiffItems:   diffItems,
		Violations:  violations,
		DiffSummary: differ.Render(diffItems, formatOpts) + "\n" + summary,
	}

	fmt := reporter.Format(reportFormat)
	if reportOutput != "" {
		if err := reporter.WriteFile(reportOutput, r, fmt); err != nil {
			return fmt.Errorf("writing report: %w", err)
		}
		cmd.Printf("Report written to %s\n", reportOutput)
		return nil
	}
	return reporter.Write(os.Stdout, r, fmt)
}
