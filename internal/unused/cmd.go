package unused

import (
	"fmt"
	"os"

	"envkit/internal/config"
	"envkit/internal/env"
)

// RunUnused is the CLI entry point for the `envkit unused` command.
// It loads the config, reads each target's env file, and prints any
// unused or undeclared variables.
func RunUnused(configPath string, targets []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	envMaps := make(map[string]map[string]string)

	for _, target := range cfg.Targets {
		if len(targets) > 0 && !contains(targets, target.Name) {
			continue
		}

		if target.EnvFile == "" {
			envMaps[target.Name] = map[string]string{}
			continue
		}

		entries, err := env.ParseFile(target.EnvFile)
		if err != nil {
			// Treat a missing env file as an empty map rather than a hard error.
			envMaps[target.Name] = map[string]string{}
			continue
		}

		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		envMaps[target.Name] = m
	}

	report := Check(cfg, envMaps)

	if !report.HasIssues() {
		fmt.Println("✓", report.Summary())
		return nil
	}

	for _, r := range report.Results {
		fmt.Fprintf(os.Stderr, "  [%s] %s — %s\n", r.Target, r.Key, r.Reason)
	}
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "✗", report.Summary())
	return fmt.Errorf("unused check failed")
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
