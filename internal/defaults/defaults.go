package defaults

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds the outcome of applying a default value to a variable.
type Result struct {
	Key     string
	OldVal  string
	NewVal  string
	Applied bool
	Reason  string
}

// Options controls how defaults are applied.
type Options struct {
	Target    string
	Keys      []string
	Overwrite bool
}

// Apply fills in missing (empty) variable values with their declared defaults
// for the given target. If Overwrite is true, existing non-empty values are
// also replaced with the declared default.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	target := findTarget(cfg, opts.Target)
	if target == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	filter := toSet(opts.Keys)
	var results []Result

	for i, v := range target.Vars {
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}
		if v.Default == "" {
			continue
		}

		result := Result{Key: v.Key, OldVal: v.Value, NewVal: v.Default}

		if v.Value == "" {
			target.Vars[i].Value = v.Default
			result.Applied = true
			result.NewVal = v.Default
			result.Reason = "was empty"
		} else if opts.Overwrite {
			target.Vars[i].Value = v.Default
			result.Applied = true
			result.Reason = "overwrite requested"
		} else {
			result.Applied = false
			result.NewVal = v.Value
			result.Reason = "already set"
		}

		results = append(results, result)
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

func toSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
