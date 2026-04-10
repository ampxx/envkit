package promote

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of a single variable promotion.
type Result struct {
	Key      string
	From     string
	To       string
	OldValue string
	NewValue string
	Skipped  bool
	Reason   string
}

// Options controls promotion behaviour.
type Options struct {
	// Keys limits promotion to specific keys; empty means all.
	Keys []string
	// Overwrite replaces an existing value in the target.
	Overwrite bool
}

// Promote copies variable values from one target to another and returns
// per-key results. The config is mutated in place; callers are responsible
// for persisting it.
func Promote(cfg *config.Config, from, to string, opts Options) ([]Result, error) {
	srcTarget := findTarget(cfg, from)
	if srcTarget == nil {
		return nil, fmt.Errorf("source target %q not found", from)
	}
	dstTarget := findTarget(cfg, to)
	if dstTarget == nil {
		return nil, fmt.Errorf("destination target %q not found", to)
	}

	filter := toSet(opts.Keys)
	var results []Result

	for _, v := range srcTarget.Vars {
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}
		res := Result{Key: v.Key, From: from, To: to, NewValue: v.Default}
		existing := findVar(dstTarget, v.Key)
		if existing != nil {
			res.OldValue = existing.Default
			if !opts.Overwrite {
				res.Skipped = true
				res.Reason = "already exists (use --overwrite to replace)"
				results = append(results, res)
				continue
			}
			existing.Default = v.Default
		} else {
			newVar := v
			dstTarget.Vars = append(dstTarget.Vars, newVar)
		}
		results = append(results, res)
	}
	return results, nil
}

func findTarget(cfg *config.Config, name string) *config.Target {
	for i := range cfg.Targets {
		if strings.EqualFold(cfg.Targets[i].Name, name) {
			return &cfg.Targets[i]
		}
	}
	return nil
}

func findVar(t *config.Target, key string) *config.VarDef {
	for i := range t.Vars {
		if t.Vars[i].Key == key {
			return &t.Vars[i]
		}
	}
	return nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
