package envimport

import (
	"fmt"
	"strings"

	"envkit/internal/config"
	"envkit/internal/env"
)

// Result describes the outcome of a single variable import.
type Result struct {
	Key     string
	Value   string
	Status  string // "imported", "skipped", "updated"
}

// Options controls how Import behaves.
type Options struct {
	Target    string
	Overwrite bool
	DryRun    bool
	Keys      []string // if non-empty, only import these keys
}

// Apply reads key=value pairs from envFile and merges them into the
// named target in cfg. It returns one Result per processed variable.
func Apply(cfg *config.Document, envFile string, opts Options) ([]Result, error) {
	entries, err := env.ParseFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("envimport: parse %q: %w", envFile, err)
	}

	filter := toSet(opts.Keys)

	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == opts.Target {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("envimport: target %q not found", opts.Target)
	}

	existing := existingKeys(target)

	var results []Result
	for _, e := range entries {
		key := strings.TrimSpace(e.Key)
		if len(filter) > 0 && !filter[key] {
			continue
		}

		_, exists := existing[key]
		switch {
		case exists && !opts.Overwrite:
			results = append(results, Result{Key: key, Value: e.Value, Status: "skipped"})
		case exists && opts.Overwrite:
			if !opts.DryRun {
				upsertVar(target, key, e.Value)
			}
			results = append(results, Result{Key: key, Value: e.Value, Status: "updated"})
		default:
			if !opts.DryRun {
				upsertVar(target, key, e.Value)
			}
			results = append(results, Result{Key: key, Value: e.Value, Status: "imported"})
		}
	}
	return results, nil
}

func existingKeys(t *config.Target) map[string]struct{} {
	m := make(map[string]struct{}, len(t.Vars))
	for _, v := range t.Vars {
		m[v.Key] = struct{}{}
	}
	return m
}

func upsertVar(t *config.Target, key, value string) {
	for i := range t.Vars {
		if t.Vars[i].Key == key {
			t.Vars[i].Value = value
			return
		}
	}
	t.Vars = append(t.Vars, config.VarDef{Key: key, Value: value})
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
