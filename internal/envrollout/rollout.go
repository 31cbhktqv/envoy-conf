package envrollout

import "sort"

// Stage represents a single deployment stage with its env vars.
type Stage struct {
	Name string
	Env  map[string]string
}

// Result holds the rollout readiness outcome for a stage transition.
type Result struct {
	From        string
	To          string
	MissingKeys []string
	ChangedKeys []string
	NewKeys     []string
	Ready       bool
}

// Plan evaluates a sequence of stages and returns per-transition results.
// Each stage is compared against the next; the final slice has len(stages)-1 entries.
func Plan(stages []Stage) []Result {
	results := make([]Result, 0, len(stages)-1)
	for i := 0; i < len(stages)-1; i++ {
		a := stages[i]
		b := stages[i+1]
		r := compare(a, b)
		results = append(results, r)
	}
	return results
}

// HasBlocker reports whether any result in the plan is not ready.
func HasBlocker(results []Result) bool {
	for _, r := range results {
		if !r.Ready {
			return true
		}
	}
	return false
}

func compare(a, b Stage) Result {
	r := Result{From: a.Name, To: b.Name}

	for k := range a.Env {
		if _, ok := b.Env[k]; !ok {
			r.MissingKeys = append(r.MissingKeys, k)
		} else if a.Env[k] != b.Env[k] {
			r.ChangedKeys = append(r.ChangedKeys, k)
		}
	}
	for k := range b.Env {
		if _, ok := a.Env[k]; !ok {
			r.NewKeys = append(r.NewKeys, k)
		}
	}

	sort.Strings(r.MissingKeys)
	sort.Strings(r.ChangedKeys)
	sort.Strings(r.NewKeys)

	r.Ready = len(r.MissingKeys) == 0
	return r
}
