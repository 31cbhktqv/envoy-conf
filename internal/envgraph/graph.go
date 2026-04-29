package envgraph

import (
	"fmt"
	"sort"
	"strings"
)

// Node represents a single environment variable and its dependencies.
type Node struct {
	Key  string
	Deps []string
}

// Graph holds the dependency relationships between env vars.
type Graph struct {
	nodes map[string]*Node
}

// New creates an empty Graph.
func New() *Graph {
	return &Graph{nodes: make(map[string]*Node)}
}

// Add registers a key with its list of dependency keys.
func (g *Graph) Add(key string, deps []string) {
	g.nodes[key] = &Node{Key: key, Deps: deps}
}

// Build derives the graph from an env map by scanning values for ${VAR} references.
func Build(env map[string]string) *Graph {
	g := New()
	for k, v := range env {
		deps := extractRefs(v)
		g.Add(k, deps)
	}
	return g
}

// Order returns a topologically sorted list of keys, or an error if a cycle exists.
func (g *Graph) Order() ([]string, error) {
	visited := make(map[string]int) // 0=unvisited,1=visiting,2=done
	var result []string

	var visit func(k string) error
	visit = func(k string) error {
		switch visited[k] {
		case 2:
			return nil
		case 1:
			return fmt.Errorf("cycle detected at %q", k)
		}
		visited[k] = 1
		if n, ok := g.nodes[k]; ok {
			for _, dep := range n.Deps {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		visited[k] = 2
		result = append(result, k)
		return nil
	}

	keys := make([]string, 0, len(g.nodes))
	for k := range g.nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if err := visit(k); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// extractRefs finds all ${VAR} or $VAR style references in a value string.
func extractRefs(value string) []string {
	var refs []string
	seen := make(map[string]bool)
	remaining := value
	for {
		start := strings.Index(remaining, "$")
		if start < 0 {
			break
		}
		remaining = remaining[start+1:]
		var name string
		if strings.HasPrefix(remaining, "{") {
			end := strings.Index(remaining, "}")
			if end < 0 {
				break
			}
			name = remaining[1:end]
			remaining = remaining[end+1:]
		} else {
			end := strings.IndexFunc(remaining, func(r rune) bool {
				return !(r == '_' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
			})
			if end < 0 {
				name = remaining
				remaining = ""
			} else {
				name = remaining[:end]
				remaining = remaining[end:]
			}
		}
		if name != "" && !seen[name] {
			seen[name] = true
			refs = append(refs, name)
		}
	}
	return refs
}
