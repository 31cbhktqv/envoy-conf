package envpin

import (
	"fmt"
	"sort"
	"time"
)

// PinnedEnv represents a point-in-time pin of environment variable values.
type PinnedEnv struct {
	Target    string            `json:"target"`
	PinnedAt  time.Time         `json:"pinned_at"`
	Variables map[string]string `json:"variables"`
}

// DriftEntry describes a single variable that has changed relative to a pin.
type DriftEntry struct {
	Key      string
	Pinned   string
	Current  string
	Status   string // "changed", "added", "removed"
}

// Pin creates a new PinnedEnv from the provided environment map.
func Pin(target string, env map[string]string) *PinnedEnv {
	vars := make(map[string]string, len(env))
	for k, v := range env {
		vars[k] = v
	}
	return &PinnedEnv{
		Target:   target,
		PinnedAt: time.Now().UTC(),
		Variables: vars,
	}
}

// Compare checks the current env against a pin and returns any drift entries.
func Compare(pin *PinnedEnv, current map[string]string) []DriftEntry {
	var entries []DriftEntry

	// Keys in pin but missing or changed in current
	for k, pinnedVal := range pin.Variables {
		curVal, ok := current[k]
		if !ok {
			entries = append(entries, DriftEntry{Key: k, Pinned: pinnedVal, Current: "", Status: "removed"})
		} else if curVal != pinnedVal {
			entries = append(entries, DriftEntry{Key: k, Pinned: pinnedVal, Current: curVal, Status: "changed"})
		}
	}

	// Keys in current but not in pin
	for k, curVal := range current {
		if _, ok := pin.Variables[k]; !ok {
			entries = append(entries, DriftEntry{Key: k, Pinned: "", Current: curVal, Status: "added"})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// Summary returns a human-readable summary line for a pin comparison.
func Summary(target string, entries []DriftEntry) string {
	if len(entries) == 0 {
		return fmt.Sprintf("[%s] pinned env matches current — no drift detected", target)
	}
	added, removed, changed := 0, 0, 0
	for _, e := range entries {
		switch e.Status {
		case "added":
			added++
		case "removed":
			removed++
		case "changed":
			changed++
		}
	}
	return fmt.Sprintf("[%s] drift detected: +%d added, -%d removed, ~%d changed", target, added, removed, changed)
}
