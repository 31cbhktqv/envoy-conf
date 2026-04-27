package envdrift

import "sort"

// DriftStatus represents the kind of drift detected for a key.
type DriftStatus string

const (
	StatusMatch   DriftStatus = "match"
	StatusAdded   DriftStatus = "added"
	StatusRemoved DriftStatus = "removed"
	StatusChanged DriftStatus = "changed"
)

// DriftEntry describes the drift state of a single environment variable.
type DriftEntry struct {
	Key      string
	Status   DriftStatus
	Baseline string
	Live     string
}

// Report holds all drift entries for a comparison.
type Report struct {
	Target  string
	Entries []DriftEntry
}

// HasDrift returns true if any entry is not a match.
func (r *Report) HasDrift() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Counts returns the number of added, removed, and changed entries.
func (r *Report) Counts() (added, removed, changed int) {
	for _, e := range r.Entries {
		switch e.Status {
		case StatusAdded:
			added++
		case StatusRemoved:
			removed++
		case StatusChanged:
			changed++
		}
	}
	return
}

// Detect compares a baseline env map against a live env map and returns a Report.
func Detect(target string, baseline, live map[string]string) Report {
	seen := make(map[string]bool)
	var entries []DriftEntry

	for k, bv := range baseline {
		seen[k] = true
		if lv, ok := live[k]; !ok {
			entries = append(entries, DriftEntry{Key: k, Status: StatusRemoved, Baseline: bv})
		} else if lv != bv {
			entries = append(entries, DriftEntry{Key: k, Status: StatusChanged, Baseline: bv, Live: lv})
		} else {
			entries = append(entries, DriftEntry{Key: k, Status: StatusMatch, Baseline: bv, Live: lv})
		}
	}

	for k, lv := range live {
		if !seen[k] {
			entries = append(entries, DriftEntry{Key: k, Status: StatusAdded, Live: lv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Report{Target: target, Entries: entries}
}
