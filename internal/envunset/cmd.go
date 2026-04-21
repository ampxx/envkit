package envunset

import (
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envkit/internal/config"
)

// RunUnset is the CLI entry-point for the `envkit unset` command.
//
// Usage:
//
	//   envkit unset KEY[,KEY...] [--target=<name>] [--dry-run] <config-file>
func RunUnset(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envkit unset <keys> <config-file> [--target=<name>] [--dry-run]")
	}

	keys := parseCSV(args[0])
	cfgPath := args[1]

	var target string
	dryRun := false
	for _, a := range args[2:] {
		if strings.HasPrefix(a, "--target=") {
			target = strings.TrimPrefix(a, "--target=")
		} else if a == "--dry-run" {
			dryRun = true
		}
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("envunset: load config: %w", err)
	}

	results, err := Apply(cfg, Options{
		Keys:   keys,
		Target: target,
		DryRun: dryRun,
	})
	if err != nil {
		return err
	}

	PrintReport(results)

	if dryRun {
		fmt.Fprintln(os.Stderr, "dry-run: no changes written")
		return nil
	}

	if err := config.Save(cfg, cfgPath); err != nil {
		return fmt.Errorf("envunset: save config: %w", err)
	}
	return nil
}

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
