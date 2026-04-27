package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"envoy-conf/internal/envloader"
	"envoy-conf/internal/envwatch"
)

var (
	watchInterval int
	watchVerbose  bool
	watchColor    bool
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <file>",
		Short: "Poll an env file for changes and stream diffs to stdout",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 5, "Poll interval in seconds")
	watchCmd.Flags().BoolVarP(&watchVerbose, "verbose", "v", false, "Show old and new values on change")
	watchCmd.Flags().BoolVar(&watchColor, "color", true, "Colorize output")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	path := args[0]

	poll := func() (map[string]string, error) {
		return envloader.LoadFile(path)
	}

	opts := envwatch.Options{
		Interval: time.Duration(watchInterval) * time.Second,
		MaxPolls: 0,
	}

	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		close(done)
	}()

	fmt.Fprintf(os.Stdout, "watching %s (interval: %ds) — press Ctrl+C to stop\n", path, watchInterval)

	changes, errs := envwatch.Watch(poll, opts, done)
	render := envwatch.RenderOptions{Color: watchColor, Verbose: watchVerbose}

	for {
		select {
		case batch, ok := <-changes:
			if !ok {
				return nil
			}
			envwatch.Render(os.Stdout, batch, render)
			fmt.Fprintln(os.Stdout, "  ["+envwatch.RenderSummary(batch)+"]")
		case err, ok := <-errs:
			if !ok {
				return nil
			}
			if err != nil {
				return fmt.Errorf("watch error: %w", err)
			}
		}
	}
}
