package envwatch

import (
	"fmt"
	"time"
)

// ChangeType describes what happened to a key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change represents a single detected change between two polls.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// PollFunc is a function that returns the current env map.
type PollFunc func() (map[string]string, error)

// Options configures the watcher.
type Options struct {
	Interval time.Duration
	MaxPolls int // 0 = unlimited
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Interval: 5 * time.Second,
		MaxPolls: 0,
	}
}

// Watch polls the provided PollFunc at the given interval and sends detected
// changes on the returned channel. The caller must close done to stop watching.
func Watch(poll PollFunc, opts Options, done <-chan struct{}) (<-chan []Change, <-chan error) {
	changes := make(chan []Change)
	errs := make(chan error, 1)

	go func() {
		defer close(changes)
		defer close(errs)

		prev, err := poll()
		if err != nil {
			errs <- fmt.Errorf("initial poll: %w", err)
			return
		}

		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()
		polls := 0

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				curr, err := poll()
				if err != nil {
					errs <- fmt.Errorf("poll: %w", err)
					return
				}
				if diff := diff(prev, curr); len(diff) > 0 {
					changes <- diff
				}
				prev = curr
				polls++
				if opts.MaxPolls > 0 && polls >= opts.MaxPolls {
					return
				}
			}
		}
	}()

	return changes, errs
}

func diff(prev, curr map[string]string) []Change {
	var out []Change
	for k, v := range curr {
		old, ok := prev[k]
		if !ok {
			out = append(out, Change{Key: k, Type: Added, NewVal: v})
		} else if old != v {
			out = append(out, Change{Key: k, Type: Changed, OldVal: old, NewVal: v})
		}
	}
	for k, v := range prev {
		if _, ok := curr[k]; !ok {
			out = append(out, Change{Key: k, Type: Removed, OldVal: v})
		}
	}
	return out
}
