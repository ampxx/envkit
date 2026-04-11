package rename

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds the outcome of a single rename operation.
type Result struct {
	Target  string
	OldKey  string
	NewKey  string
	Renamed bool
	Reason  string
}

// Options controls rename behaviour.
type Options struct {
	Target    string // empty means all targets
	Keys      map[string]string // oldKey -> newKey
	Overwrite bool // if true, overwrite existing newKey
}

// Apply renames variable keys in the config according to opts.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	if len(opts.Keys) == 0 {
		return nil, fmt.Errorf("rename: no key mappings provided")
	}

	var results []Result

	for i, target := range cfg.Targets {
		if opts.Target != "" && target.Name != opts.Target {
			continue
		}

		for oldKey, newKey := range opts.Keys {
			oldIdx := findVar(target.Vars, oldKey)
			if oldIdx < 0 {
				results = append(results, Result{
					Target: target.Name,
					OldKey: oldKey,
					NewKey: newKey,
					Renamed: false,
					Reason: "key not found",
				})
				continue
			}

			newIdx := findVar(target.Vars, newKey)
			if newIdx >= 0 && !opts.Overwrite {
				results = append(results, Result{
					Target: target.Name,
					OldKey: oldKey,
					NewKey: newKey,
					Renamed: false,
					Reason: "destination key already exists",
				})
				continue
			}

			// Remove destination if overwriting
			if newIdx >= 0 {
				cfg.Targets[i].Vars = removeVar(cfg.Targets[i].Vars, newIdx)
				// recalculate oldIdx after removal
				oldIdx = findVar(cfg.Targets[i].Vars, oldKey)
			}

			cfg.Targets[i].Vars[oldIdx].Key = newKey
			results = append(results, Result{
				Target:  target.Name,
				OldKey:  oldKey,
				NewKey:  newKey,
				Renamed: true,
			})
		}
	}

	return results, nil
}

func findVar(vars []config.VarDef, key string) int {
	for i, v := range vars {
		if v.Key == key {
			return i
		}
	}
	return -1
}

func removeVar(vars []config.VarDef, idx int) []config.VarDef {
	return append(vars[:idx], vars[idx+1:]...)
}
