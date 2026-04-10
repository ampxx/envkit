package compare

import (
	"fmt"
	"os"

	"github.com/user/envkit/internal/config"
)

// RunCompare loads the config file and prints a diff between two named targets.
func RunCompare(cfgPath, targetA, targetB string) error {
	if targetA == "" || targetB == "" {
		return fmt.Errorf("two target names are required")
	}
	if targetA == targetB {
		return fmt.Errorf("targets must be different")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	diff, err := Targets(cfg, targetA, targetB)
	if err != nil {
		return err
	}

	PrintReport(diff)

	if len(diff.Differing) > 0 || len(diff.OnlyInA) > 0 || len(diff.OnlyInB) > 0 {
		os.Exit(1)
	}
	return nil
}
