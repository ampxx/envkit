package envset

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunSet is the CLI entry-point for the `envkit set` command.
// args is a list of KEY=VALUE strings.
func RunSet(cfgPath, target string, args []string, overwrite, dryRun bool) error {
	if len(args) == 0 {
		return fmt.Errorf("provide at least one KEY=VALUE pair")
	}

	pairs, err := parsePairs(args)
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	opts := Options{
		Target:    target,
		Overwrite: overwrite,
		DryRun:    dryRun,
	}

	results, err := Apply(cfg, pairs, opts)
	if err != nil {
		return err
	}

	PrintReport(results)

	if dryRun {
		fmt.Fprintln(os.Stderr, "dry-run: no changes written")
		return nil
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	return nil
}

func parsePairs(args []string) (map[string]string, error) {
	pairs := make(map[string]string, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid pair %q: expected KEY=VALUE", arg)
		}
		pairs[parts[0]] = parts[1]
	}
	return pairs, nil
}
