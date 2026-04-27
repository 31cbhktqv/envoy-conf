package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envexport"
	"envoy-conf/internal/envloader"
	"envoy-conf/internal/masker"
)

var (
	exportFormat     string
	exportMaskKeys   []string
	exportAutoMask   bool
	exportOutputFile string
)

func init() {
	exportCmd := &cobra.Command{
		Use:   "export <env-file>",
		Short: "Export env vars to dotenv, shell export, JSON, or YAML",
		Args:  cobra.ExactArgs(1),
		RunE:  runExport,
	}

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv|export|json|yaml")
	exportCmd.Flags().StringSliceVar(&exportMaskKeys, "mask", nil, "Keys to mask in output")
	exportCmd.Flags().BoolVar(&exportAutoMask, "auto-mask", false, "Auto-detect and mask sensitive keys")
	exportCmd.Flags().StringVarP(&exportOutputFile, "output", "o", "", "Write output to file instead of stdout")

	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	env, err := envloader.LoadFile(args[0])
	if err != nil {
		return fmt.Errorf("loading env file: %w", err)
	}

	maskedKeys := map[string]bool{}
	for _, k := range exportMaskKeys {
		maskedKeys[strings.TrimSpace(k)] = true
	}

	if exportAutoMask {
		m, err := masker.New(nil)
		if err != nil {
			return fmt.Errorf("creating masker: %w", err)
		}
		for k := range env {
			if m.IsSensitive(k) {
				maskedKeys[k] = true
			}
		}
	}

	opts := envexport.DefaultOptions()
	opts.Format = envexport.Format(exportFormat)
	if len(maskedKeys) > 0 {
		opts.Masked = maskedKeys
	}

	out, err := envexport.Export(env, opts)
	if err != nil {
		return err
	}

	if exportOutputFile != "" {
		return os.WriteFile(exportOutputFile, []byte(out), 0o644)
	}

	fmt.Print(out)
	return nil
}
