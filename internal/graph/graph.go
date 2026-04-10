package graph

import (
	"fmt"
	"sort"

	"envkit/internal/config"
)

// Node represents a variable and its dependencies.
type Node struct {
	Key  string
	Deps []string
}

// Graph holds the dependency graph for environment variables.
type Graph struct {
	Nodes map[string]*Node
}

// Build constructs a dependency graph from a config by scanning
// variable default values for ${VAR} references.
func Build(cfg *config.Config, target string) (*Graph, error) {
	g := &Graph{Nodes: make(map[string]*Node)}

	vars := cfg.Vars
	for _, t := range cfg.Targets {
		if t.Name == target {
			for k, v := range t.Vars {
				vars[k] = v
			}
			break
		}
	}

	for key, def := range vars {
		node := &Node{Key: key}
		node.Deps = extractRefs(def.Default)
		g.Nodes[key] = node
	}

	if err := g.detectCycles(); err != nil {
		return nil, err
	}
	return g, nil
}

// Order returns a topologically sorted list of keys.
func (g *Graph) Order() ([]string, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var result []string

	keys := make([]string, 0, len(g.Nodes))
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var visit func(k string) error
	visit = func(k string) error {
		if temp[k] {
			return fmt.Errorf("cycle detected at %q", k)
		}
		if visited[k] {
			return nil
		}
		temp[k] = true
		if node, ok := g.Nodes[k]; ok {
			for _, dep := range node.Deps {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		temp[k] = false
		visited[k] = true
		result = append(result, k)
		return nil
	}

	for _, k := range keys {
		if err := visit(k); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (g *Graph) detectCycles() error {
	_, err := g.Order()
	return err
}

// extractRefs parses ${VAR} patterns from a string.
func extractRefs(s string) []string {
	var refs []string
	for i := 0; i < len(s); i++ {
		if i+2 < len(s) && s[i] == '$' && s[i+1] == '{' {
			end := i + 2
			for end < len(s) && s[end] != '}' {
				end++
			}
			if end < len(s) {
				refs = append(refs, s[i+2:end])
				i = end
			}
		}
	}
	return refs
}
