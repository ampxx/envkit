package graph

import (
	"fmt"
	"os"

	"envkit/internal/config"
)

// RunGraph loads the config and prints the dependency graph for a given target.
func RunGraph(configPath, target string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	g, err := Build(cfg, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	if target != "" {
		fmt.Printf("Target: %s\n", target)
	} else {
		fmt.Println("Target: (global)")
	}

	PrintGraph(g)
	fmt.Println()
	fmt.Println(Summary(g))
	return nil
}
