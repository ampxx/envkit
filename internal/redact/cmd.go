package redact

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/yourusername/envkit/internal/config"
	"github.com/yourusername/envkit/internal/env"
)

// RunRedact loads an env file and prints its values with sensitive ones masked.
// configPath: path to envkit config, envFile: .env file to read, target: target name.
func RunRedact(configPath, envFile, targetName string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	parsed, err := env.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}

	vars := env.ToMap(parsed)

	// If a target is specified, restrict to keys declared in that target.
	if targetName != "" {
		filtered, err := filterByTarget(cfg, vars, targetName)
		if err != nil {
			return err
		}
		vars = filtered
	}

	results := Apply(vars, nil)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE\tREDACTED")
	fmt.Fprintln(w, "---\t-----\t--------")
	for _, r := range results {
		redactedMark := ""
		if r.Redacted {
			redactedMark = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Key, r.Value, redactedMark)
	}
	return w.Flush()
}

// filterByTarget returns a subset of vars containing only keys declared in the
// named target. It returns an error if the target is not found in the config.
func filterByTarget(cfg *config.Config, vars map[string]string, targetName string) (map[string]string, error) {
	for _, t := range cfg.Targets {
		if t.Name != targetName {
			continue
		}
		filtered := make(map[string]string, len(t.Vars))
		for _, v := range t.Vars {
			if val, ok := vars[v.Key]; ok {
				filtered[v.Key] = val
			}
		}
		return filtered, nil
	}
	return nil, fmt.Errorf("target %q not found in config", targetName)
}
