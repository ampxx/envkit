package schema

import (
	"fmt"
	"os"

	"github.com/envkit/envkit/internal/config"
)

// RunSchema loads the config at cfgPath, validates variable schema for the
// given target, and prints any issues to stdout. Exits with code 1 if issues
// are found.
func RunSchema(cfgPath, target string) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	issues := Validate(cfg, target)
	if len(issues) == 0 {
		fmt.Printf("schema OK — no issues found in target %q\n", target)
		return
	}

	fmt.Printf("schema issues found in target %q:\n\n", target)
	for _, iss := range issues {
		fmt.Printf("  %-30s %s\n", iss.Key, iss.Message)
	}
	fmt.Printf("\n%d issue(s) detected.\n", len(issues))
	os.Exit(1)
}
