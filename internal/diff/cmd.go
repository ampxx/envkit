package diff

import (
	"fmt"
	"os"

	"github.com/envkit/envkit/internal/config"
)

// RunDiff loads two named targets from the config file and prints their diff.
func RunDiff(cfgPath, targetA, targetB string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	a, err := config.FindTarget(cfg, targetA)
	if err != nil {
		return fmt.Errorf("target %q: %w", targetA, err)
	}

	b, err := config.FindTarget(cfg, targetB)
	if err != nil {
		return fmt.Errorf("target %q: %w", targetB, err)
	}

	result := Compare(a.Vars, b.Vars)
	output := Format(result, targetA, targetB)

	fmt.Fprint(os.Stdout, output)

	if len(result.Differing) > 0 || len(result.OnlyInA) > 0 || len(result.OnlyInB) > 0 {
		return fmt.Errorf("%d key(s) differ between %s and %s",
			len(result.Differing)+len(result.OnlyInA)+len(result.OnlyInB),
			targetA, targetB)
	}

	return nil
}
