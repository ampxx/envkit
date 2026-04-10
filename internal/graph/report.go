package graph

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PrintGraph writes a human-readable dependency graph to stdout.
func PrintGraph(g *Graph) {
	printGraphTo(os.Stdout, g)
}

func printGraphTo(w io.Writer, g *Graph) {
	order, err := g.Order()
	if err != nil {
		fmt.Fprintf(w, "error computing order: %v\n", err)
		return
	}

	fmt.Fprintln(w, "Dependency Graph (topological order):")
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, key := range order {
		node := g.Nodes[key]
		if len(node.Deps) == 0 {
			fmt.Fprintf(w, "  %s\n", key)
		} else {
			fmt.Fprintf(w, "  %s -> depends on: [%s]\n", key, strings.Join(node.Deps, ", "))
		}
	}
}

// Summary returns a one-line summary of the graph.
func Summary(g *Graph) string {
	leaves := 0
	for _, n := range g.Nodes {
		if len(n.Deps) == 0 {
			leaves++
		}
	}
	return fmt.Sprintf("%d variables, %d with dependencies, %d independent",
		len(g.Nodes), len(g.Nodes)-leaves, leaves)
}
