package inherit

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds the outcome of a single variable inheritance.
type Result struct {
	Key      string
	FromTarget string
	ToTarget   string
	Value    string
	Skipped  bool
	Reason   string
}

// Options controls how inheritance is applied.
type Options struct {
	Keys      []string // if empty, all vars are considered
	Overwrite bool
}

// Apply copies variables from srcTarget into dstTarget, inheriting any
// vars that exist in src but are absent (or empty) in dst.
func Apply(cfg *config.Document, srcTarget, dstTarget string, opts Options) ([]Result, error) {
	src := findTarget(cfg, srcTarget)
	if src == nil {
		return nil, fmt.Errorf("source target %q not found", srcTarget)
	}
	dst := findTarget(cfg, dstTarget)
	if dst == nil {
		return nil, fmt.Errorf("destination target %q not found", dstTarget)
	}

	filter := toSet(opts.Keys)
	dstKeys := existingKeys(dst)

	var results []Result
	for _, v := range src.Vars {
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}
		_, exists := dstKeys[v.Key]
		if exists && !opts.Overwrite {
			results = append(results, Result{
				Key: v.Key, FromTarget: srcTarget, ToTarget: dstTarget,
				Value: v.Value, Skipped: true, Reason: "already exists",
			})
			continue
		}
		upsertVar(dst, v)
		results = append(results, Result{
			Key: v.Key, FromTarget: srcTarget, ToTarget: dstTarget,
			Value: v.Value, Skipped: false,
		})
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

func upsertVar(t *config.Target, v config.VarDef) {
	for i := range t.Vars {
		if t.Vars[i].Key == v.Key {
			t.Vars[i] = v
			return
		}
	}
	t.Vars = append(t.Vars, v)
}

func toSet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
