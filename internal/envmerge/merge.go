package envmerge

import "fmt"

// Strategy defines how conflicting keys are handled during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first source that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last source that defines the key.
	StrategyLast
	// StrategyStrict returns an error if any key appears in more than one source.
	StrategyStrict
)

// Result holds the merged environment map along with metadata about conflicts.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Conflict records a key that appeared in multiple sources.
type Conflict struct {
	Key    string
	Values []string // one entry per source, in order
}

// Merge combines multiple environment maps according to the given strategy.
// Sources are applied in order; index 0 is considered the "base" source.
func Merge(strategy Strategy, sources ...map[string]string) (Result, error) {
	merged := make(map[string]string)
	conflictMap := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			if existing, found := merged[k]; found {
				conflictMap[k] = append(conflictMap[k], existing, v)
				switch strategy {
				case StrategyStrict:
					return Result{}, fmt.Errorf("merge conflict: key %q defined in multiple sources", k)
				case StrategyFirst:
					// keep existing — do nothing
				case StrategyLast:
					merged[k] = v
				}
			} else {
				merged[k] = v
			}
		}
	}

	var conflicts []Conflict
	for k, vals := range conflictMap {
		conflicts = append(conflicts, Conflict{Key: k, Values: vals})
	}

	return Result{Env: merged, Conflicts: conflicts}, nil
}
