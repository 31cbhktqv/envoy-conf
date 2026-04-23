package differ

import "sort"

// DiffResult holds the comparison between two env maps.
type DiffResult struct {
	OnlyInA   map[string]string
	OnlyInB   map[string]string
	Changed   map[string][2]string // key -> [valueA, valueB]
	Unchanged map[string]string
}

// Diff compares two environment variable maps (a and b) and returns a DiffResult.
func Diff(a, b map[string]string) DiffResult {
	result := DiffResult{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				result.Unchanged[k] = va
			} else {
				result.Changed[k] = [2]string{va, vb}
			}
		} else {
			result.OnlyInA[k] = va
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			result.OnlyInB[k] = vb
		}
	}

	return result
}

// HasDifferences returns true if there are any additions, removals, or changes.
func (d DiffResult) HasDifferences() bool {
	return len(d.OnlyInA) > 0 || len(d.OnlyInB) > 0 || len(d.Changed) > 0
}

// SortedChangedKeys returns changed keys in sorted order.
func (d DiffResult) SortedChangedKeys() []string {
	keys := make([]string, 0, len(d.Changed))
	for k := range d.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedOnlyInAKeys returns OnlyInA keys in sorted order.
func (d DiffResult) SortedOnlyInAKeys() []string {
	keys := make([]string, 0, len(d.OnlyInA))
	for k := range d.OnlyInA {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedOnlyInBKeys returns OnlyInB keys in sorted order.
func (d DiffResult) SortedOnlyInBKeys() []string {
	keys := make([]string, 0, len(d.OnlyInB))
	for k := range d.OnlyInB {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
