package envcanary

import "fmt"

// Status represents the result of a canary check for a single key.
type Status int

const (
	StatusOK      Status = iota
	StatusWarning        // value changed but within threshold
	StatusCritical       // value changed beyond threshold or key missing
)

// Result holds the canary check outcome for one environment variable.
type Result struct {
	Key      string
	Status   Status
	Baseline string
	Current  string
	Message  string
}

// Options controls canary check behaviour.
type Options struct {
	// RequiredKeys are keys that must be present in current env.
	RequiredKeys []string
	// WatchKeys are keys whose value changes should be flagged.
	WatchKeys []string
	// AllowMissing allows keys absent from current without marking critical.
	AllowMissing bool
}

// Check compares a baseline environment against a current environment and
// returns per-key canary results.
func Check(baseline, current map[string]string, opts Options) []Result {
	var results []Result

	for _, key := range opts.RequiredKeys {
		curVal, ok := current[key]
		if !ok {
			status := StatusCritical
			if opts.AllowMissing {
				status = StatusWarning
			}
			results = append(results, Result{
				Key:      key,
				Status:   status,
				Baseline: baseline[key],
				Current:  "",
				Message:  fmt.Sprintf("required key %q missing from current environment", key),
			})
			continue
		}
		results = append(results, Result{
			Key:     key,
			Status:  StatusOK,
			Current: curVal,
			Message: "present",
		})
	}

	for _, key := range opts.WatchKeys {
		baseVal := baseline[key]
		curVal, ok := current[key]
		if !ok {
			results = append(results, Result{
				Key:      key,
				Status:   StatusWarning,
				Baseline: baseVal,
				Current:  "",
				Message:  fmt.Sprintf("watched key %q absent in current environment", key),
			})
			continue
		}
		if baseVal != curVal {
			results = append(results, Result{
				Key:      key,
				Status:   StatusWarning,
				Baseline: baseVal,
				Current:  curVal,
				Message:  fmt.Sprintf("watched key %q value changed", key),
			})
			continue
		}
		results = append(results, Result{
			Key:     key,
			Status:  StatusOK,
			Current: curVal,
			Message: "unchanged",
		})
	}

	return results
}

// HasCritical returns true if any result has StatusCritical.
func HasCritical(results []Result) bool {
	for _, r := range results {
		if r.Status == StatusCritical {
			return true
		}
	}
	return false
}
