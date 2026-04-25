package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/envtemplate"
)

var (
	templateEnvFile  string
	templateLookupFile string
	templateStrict   bool
	templateFallback string
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Expand variable placeholders in an env file",
		Long: `Reads an env file and expands $VAR or ${VAR} placeholders using a
second lookup env file and/or OS environment variables.`,
		RunE: runTemplate,
	}

	templateCmd.Flags().StringVarP(&templateEnvFile, "file", "f", "", "env file with placeholders (required)")
	templateCmd.Flags().StringVarP(&templateLookupFile, "lookup", "l", "", "env file used as variable source (optional)")
	templateCmd.Flags().BoolVar(&templateStrict, "strict", false, "fail on unresolved variables")
	templateCmd.Flags().StringVar(&templateFallback, "fallback", "", "value used for unresolved variables (non-strict mode)")
	_ = templateCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, _ []string) error {
	env, err := envloader.LoadFile(templateEnvFile)
	if err != nil {
		return fmt.Errorf("loading env file: %w", err)
	}

	lookup := map[string]string{}
	if templateLookupFile != "" {
		lookup, err = envloader.LoadFile(templateLookupFile)
		if err != nil {
			return fmt.Errorf("loading lookup file: %w", err)
		}
	}

	opts := envtemplate.Options{
		Strict:   templateStrict,
		Fallback: templateFallback,
	}

	result, err := envtemplate.Expand(env, lookup, opts)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, result[k])
	}
	return nil
}
