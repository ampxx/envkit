package envunset

import (
	"fmt"

	"github.com/your-org/envkit/internal/config"
)

// Result holds the outcome of an unset operation for a single variable.
type Result struct {
	Target string
	Key    string
	Found  bool
}

// Options controls how Apply behaves.
type Options struct {
	// Keys is the explicit list of keys to remove. If empty, nothing is removed.
	Keys []string
	// Target restricts the operation to a single deployment target.
	// An empty string means all targets.
	Target string
	// DryRun reports what would be removed without mutating the config.
	DryRun bool
}

// Apply removes the specified keys from matching targets in cfg.
// It returns a Result for every (target, key) pair that was considered.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	if len(opts.Keys) == 0 {
		return nil, fmt.Errorf("envunset: at least one key must be specified")
	}
	wanted := toSet(opts.Keys)
	var results []Result

	for i := range cfg.Targets {
		t := &cfg.Targets[i]
		if opts.Target != "" && t.Name != opts.Target {
			continue
		}
		var kept []config.VarDef
		for _, v := range t.Vars {
			if wanted[v.Key] {
				results = append(results, Result{Target: t.Name, Key: v.Key, Found: true})
				if opts.DryRun {
					kept = append(kept, v)
				}
				// not appended to kept when not dry-run → effectively removed
			} else {
				kept = append(kept, v)
			}
		}
		// Record keys that were not found in this target.
		for k := range wanted {
			if !containsKey(t.Vars, k) {
				results = append(results, Result{Target: t.Name, Key: k, Found: false})
			}
		}
		if !opts.DryRun {
			t.Vars = kept
		}
	}
	return results, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func containsKey(vars []config.VarDef, key string) bool {
	for _, v := range vars {
		if v.Key == key {
			return true
		}
	}
	return false
}
