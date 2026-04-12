package extract

import (
	"fmt"
	"regexp"

	"envkit/internal/config"
)

// Result holds the outcome of extracting a single variable.
type Result struct {
	Key     string
	Value   string
	Target  string
	Matched bool
}

// Options controls extraction behaviour.
type Options struct {
	Target  string
	Keys    []string // explicit keys; empty means all
	Pattern string   // regex pattern to match key names
}

// Apply extracts variables from the given target, optionally filtered by key
// list or regex pattern.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == opts.Target {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern: %w", err)
		}
	}

	keySet := toSet(opts.Keys)

	var results []Result
	for _, v := range target.Vars {
		if len(keySet) > 0 && !keySet[v.Key] {
			continue
		}
		if re != nil && !re.MatchString(v.Key) {
			continue
		}
		results = append(results, Result{
			Key:     v.Key,
			Value:   v.Default,
			Target:  opts.Target,
			Matched: true,
		})
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
