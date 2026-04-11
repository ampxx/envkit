package scope

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// Result holds the scoped variable set for a given target.
type Result struct {
	Target string
	Vars    map[string]string
}

// Options controls how scoping is applied.
type Options struct {
	Target  string
	Keys    []string // if non-empty, only include these keys
	Exclude []string // keys to exclude
}

// Apply returns the resolved variable map for the target, filtered by Options.
func Apply(cfg *config.Document, opts Options) (Result, error) {
	target := findTarget(cfg, opts.Target)
	if target == nil {
		return Result{}, fmt.Errorf("target %q not found", opts.Target)
	}

	include := toSet(opts.Keys)
	exclude := toSet(opts.Exclude)

	vars := make(map[string]string)
	for _, v := range target.Vars {
		if len(include) > 0 && !include[v.Key] {
			continue
		}
		if exclude[v.Key] {
			continue
		}
		val := v.Default
		if v.Value != "" {
			val = v.Value
		}
		vars[v.Key] = val
	}

	return Result{Target: opts.Target, Vars: vars}, nil
}

// Prefix returns a new Result with all keys prefixed by the given string.
func Prefix(r Result, prefix string) Result {
	out := make(map[string]string, len(r.Vars))
	for k, v := range r.Vars {
		out[strings.ToUpper(prefix)+"_"+k] = v
	}
	return Result{Target: r.Target, Vars: out}
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
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
