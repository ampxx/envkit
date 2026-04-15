package envdiff

import (
	"fmt"
	"os"

	"github.com/your-org/envkit/internal/config"
)

// RunEnvDiff is the entry point for the `envkit envdiff` sub-command.
// Usage: envkit envdiff <config> <targetA> <targetB> [--all]
func RunEnvDiff(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envdiff <config-file> <targetA> <targetB> [--all]")
	}

	cfgPath := args[0]
	targetA := args[1]
	targetB := args[2]
	showAll := len(args) >= 4 && args[3] == "--all"

	doc, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("envdiff: loading config: %w", err)
	}

	report, err := Compare(doc, targetA, targetB)
	if err != nil {
		return fmt.Errorf("envdiff: %w", err)
	}

	PrintReport(report, showAll)

	// Exit with non-zero if there are any differences.
	for _, r := range report.Results {
		if r.Kind != KindIdentical {
			os.Exit(1)
		}
	}
	return nil
}
