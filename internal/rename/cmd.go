package rename

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// RunRename is the CLI entry point for the rename command.
// pairs is a slice of "OLD=NEW" strings.
func RunRename(cfgPath, target string, pairs []string, overwrite bool) error {
	if len(pairs) == 0 {
		return fmt.Errorf("rename: provide at least one OLD=NEW pair")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("rename: load config: %w", err)
	}

	keys, err := parsePairs(pairs)
	if err != nil {
		return err
	}

	results, err := Apply(cfg, Options{
		Target:    target,
		Keys:      keys,
		Overwrite: overwrite,
	})
	if err != nil {
		return err
	}

	printResults(results)

	if err := config.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("rename: save config: %w", err)
	}
	return nil
}

func parsePairs(pairs []string) (map[string]string, error) {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("rename: invalid pair %q, expected OLD=NEW", p)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}

func printResults(results []Result) {
	renamed, skipped := 0, 0
	for _, r := range results {
		if r.Renamed {
			fmt.Printf("  ✔ [%s] %s → %s\n", r.Target, r.OldKey, r.NewKey)
			renamed++
		} else {
			fmt.Printf("  ✗ [%s] %s → %s (%s)\n", r.Target, r.OldKey, r.NewKey, r.Reason)
			skipped++
		}
	}
	fmt.Printf("\nRenamed: %d  Skipped: %d\n", renamed, skipped)
}
