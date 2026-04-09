package lint

import (
	"fmt"
	"os"

	"github.com/envkit/envkit/internal/config"
)

// RunLint loads the config at the given path, runs lint checks,
// prints a report, and exits with a non-zero code if errors are found.
func RunLint(cfgPath string) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	result := Run(cfg)

	if len(result.Issues) == 0 {
		fmt.Println("✔ No lint issues found.")
		return
	}

	errCount, warnCount := 0, 0
	for _, issue := range result.Issues {
		switch issue.Level {
		case "error":
			fmt.Printf("  [ERROR] %s: %s\n", issue.Key, issue.Message)
			errCount++
		case "warn":
			fmt.Printf("  [WARN]  %s: %s\n", issue.Key, issue.Message)
			warnCount++
		}
	}

	fmt.Printf("\nLint summary: %d error(s), %d warning(s)\n", errCount, warnCount)

	if result.HasErrors() {
		os.Exit(1)
	}
}
