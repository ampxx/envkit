package coerce

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
	"envkit/internal/env"
)

// RunCoerce is the CLI entry-point for the coerce command.
func RunCoerce(cfgPath, envFile, target string, keys []string, dryRun bool) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	entries, err := env.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}
	environMap := env.ToMap(entries)

	results, err := Apply(cfg, environMap, Options{
		Target: target,
		Keys:   keys,
		DryRun: dryRun,
	})
	if err != nil {
		return err
	}

	PrintReport(results)
	fmt.Fprintln(os.Stdout, Summary(results))

	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: no changes written)")
		return nil
	}

	// Write updated values back to the env file.
	var lines []string
	for _, e := range entries {
		if v, ok := environMap[e.Key]; ok {
			lines = append(lines, e.Key+"="+v)
		} else {
			lines = append(lines, e.Key+"="+e.Value)
		}
	}
	return os.WriteFile(envFile, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
