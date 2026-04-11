package trim

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// RunTrim loads the config at cfgPath, trims variable values in targetName,
// and writes the updated config back to disk.
func RunTrim(cfgPath, targetName string, trimSpace, trimQuotes bool, keys []string, dryRun bool) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	opts := Options{
		TrimSpace:  trimSpace,
		TrimQuotes: trimQuotes,
		Keys:       keys,
	}

	results, err := Apply(cfg, targetName, opts)
	if err != nil {
		return err
	}

	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
			fmt.Printf("  ~ %-20s %q -> %q\n", r.Key, r.Original, r.Trimmed)
		}
	}

	if changed == 0 {
		fmt.Println("No values needed trimming.")
		return nil
	}

	fmt.Printf("\n%d value(s) trimmed", changed)
	if dryRun {
		fmt.Println(" (dry-run, no changes written).")
		return nil
	}
	fmt.Println(".")

	if err := config.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

// parseKeys splits a comma-separated key list into a slice.
func parseKeys(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if k := strings.TrimSpace(p); k != "" {
			out = append(out, k)
		}
	}
	return out
}
