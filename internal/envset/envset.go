package envset

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds the outcome of a single set operation.
type Result struct {
	Key     string
	Value   string
	Target  string
	Created bool // true if key was new, false if updated
	Skipped bool // true if key existed and overwrite was false
}

// Options controls how Apply behaves.
type Options struct {
	Target    string
	Overwrite bool
	DryRun    bool
}

// Apply sets one or more key=value pairs in the given target within cfg.
// Returns a slice of Results describing what happened for each pair.
func Apply(cfg *config.Document, pairs map[string]string, opts Options) ([]Result, error) {
	tgt := findTarget(cfg, opts.Target)
	if tgt == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	existing := existingKeys(tgt)
	var results []Result

	for k, v := range pairs {
		r := Result{Key: k, Value: v, Target: opts.Target}

		if _, found := existing[k]; found {
			if !opts.Overwrite {
				r.Skipped = true
				results = append(results, r)
				continue
			}
			r.Created = false
		} else {
			r.Created = true
		}

		if !opts.DryRun {
			upsertVar(tgt, k, v)
		}
		results = append(results, r)
	}

	return results, nil
}

func findTarget(cfg *config.Document, name string) *config.Target {
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == name {
			return &cfg.Targets[i]
		}
	}
	return nil
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
