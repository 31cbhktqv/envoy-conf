package envpromote

import "fmt"

// Stage represents a named deployment stage (e.g. staging, production).
type Stage struct {
	Name string
	Env  map[string]string
}

// PromoteResult holds the outcome of a promotion check between two stages.
type PromoteResult struct {
	From        string
	To          string
	MissingKeys []string // keys present in From but absent in To
	NewKeys     []string // keys present in To but absent in From
	ChangedKeys []string // keys present in both but with different values
	Ready       bool
}

// Promote compares two stages and determines whether the target stage is
// ready to receive the promotion. A promotion is considered ready when there
// are no missing keys (keys that exist in the source but are absent in the
// destination).
func Promote(from, to Stage) PromoteResult {
	result := PromoteResult{
		From:  from.Name,
		To:    to.Name,
		Ready: true,
	}

	for k, vFrom := range from.Env {
		vTo, exists := to.Env[k]
		if !exists {
			result.MissingKeys = append(result.MissingKeys, k)
			result.Ready = false
		} else if vFrom != vTo {
			result.ChangedKeys = append(result.ChangedKeys, k)
		}
	}

	for k := range to.Env {
		if _, exists := from.Env[k]; !exists {
			result.NewKeys = append(result.NewKeys, k)
		}
	}

	return result
}

// Summary returns a human-readable one-line summary of the promotion result.
func Summary(r PromoteResult) string {
	if r.Ready {
		return fmt.Sprintf("promotion %s → %s: ready (%d changed, %d new)",
			r.From, r.To, len(r.ChangedKeys), len(r.NewKeys))
	}
	return fmt.Sprintf("promotion %s → %s: NOT READY (%d missing, %d changed, %d new)",
		r.From, r.To, len(r.MissingKeys), len(r.ChangedKeys), len(r.NewKeys))
}
