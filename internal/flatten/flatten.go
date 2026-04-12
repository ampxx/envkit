package flatten

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of flattening a single variable.
type Result struct {
	Key    string
	Value  string
	Source string // "global" or target name
}

// Options controls flatten behaviour.
type Options struct {
	Target    string
	Prefix    string
	Separator string // default "_"
	Keys      []string
}

// Apply flattens all resolved variables for a target into a single key=value
// map, optionally adding a prefix and filtering by key set.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	tgt := findTarget(cfg, opts.Target)
	if opts.Target != "" && tgt == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	allow := toSet(opts.Keys)

	// Start from globals.
	resolved := map[string]Result{}
	for _, v := range cfg.Vars {
		if len(allow) > 0 && !allow[v.Key] {
			continue
		}
		resolved[v.Key] = Result{Key: v.Key, Value: v.Default, Source: "global"}
	}

	// Override with target vars.
	if tgt != nil {
		for _, v := range tgt.Vars {
			if len(allow) > 0 && !allow[v.Key] {
				continue
			}
			resolved[v.Key] = Result{Key: v.Key, Value: v.Default, Source: tgt.Name}
		}
	}

	// Apply prefix.
	var out []Result
	for _, r := range resolved {
		if opts.Prefix != "" {
			r.Key = strings.ToUpper(opts.Prefix) + sep + r.Key
		}
		out = append(out, r)
	}

	// Stable sort by key.
	sortResults(out)
	return out, nil
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

func sortResults(rs []Result) {
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && rs[j].Key < rs[j-1].Key; j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}
