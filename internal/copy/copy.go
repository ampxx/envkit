package copy

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds the outcome of a single copy operation.
type Result struct {
	Key     string
	From    string
	To      string
	Skipped bool
	Reason  string
}

// Options controls copy behaviour.
type Options struct {
	Keys      []string // if empty, copy all keys
	Overwrite bool
}

// Copy copies variable definitions from one target to another within cfg.
// It returns the list of results and a modified copy of the config.
func Copy(cfg config.Document, from, to string, opts Options) ([]Result, config.Document, error) {
	src := findTarget(cfg, from)
	if src == nil {
		return nil, cfg, fmt.Errorf("source target %q not found", from)
	}
	dst := findTarget(cfg, to)
	if dst == nil {
		return nil, cfg, fmt.Errorf("destination target %q not found", to)
	}

	filter := toSet(opts.Keys)
	dstKeys := existingKeys(*dst)

	var results []Result
	for _, v := range src.Vars {
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}
		if _, exists := dstKeys[v.Key]; exists && !opts.Overwrite {
			results = append(results, Result{Key: v.Key, From: from, To: to, Skipped: true, Reason: "already exists"})
			continue
		}
		dst = upsertVar(dst, v)
		results = append(results, Result{Key: v.Key, From: from, To: to})
	}

	cfg = replaceTarget(cfg, *dst)
	return results, cfg, nil
}

func findTarget(cfg config.Document, name string) *config.Target {
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == name {
			return &cfg.Targets[i]
		}
	}
	return nil
}

func existingKeys(t config.Target) map[string]struct{} {
	m := make(map[string]struct{}, len(t.Vars))
	for _, v := range t.Vars {
		m[v.Key] = struct{}{}
	}
	return m
}

func upsertVar(t *config.Target, v config.VarDef) *config.Target {
	for i, existing := range t.Vars {
		if existing.Key == v.Key {
			t.Vars[i] = v
			return t
		}
	}
	t.Vars = append(t.Vars, v)
	return t
}

func replaceTarget(cfg config.Document, t config.Target) config.Document {
	for i, existing := range cfg.Targets {
		if existing.Name == t.Name {
			cfg.Targets[i] = t
			return cfg
		}
	}
	return cfg
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
