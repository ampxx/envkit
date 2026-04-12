package filter

import (
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of a filter operation for a single variable.
type Result struct {
	Key     string
	Value   string
	Matched bool
}

// Options controls how filtering is applied.
type Options struct {
	Target  string
	Keys    []string // explicit key allowlist (empty = all)
	Tags    []string // match vars that have ALL of these tags
	Prefix  string   // key must start with this prefix
	Pattern string   // substring match against key
}

// Apply filters variables for the given target according to opts.
// Returns only the variables that satisfy every non-empty criterion.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == opts.Target {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("filter: target %q not found", opts.Target)
	}

	keySet := toSet(opts.Keys)
	tagSet := toSet(opts.Tags)

	var results []Result
	for _, v := range target.Vars {
		if len(keySet) > 0 && !keySet[v.Key] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(v.Key, opts.Prefix) {
			continue
		}
		if opts.Pattern != "" && !strings.Contains(v.Key, opts.Pattern) {
			continue
		}
		if len(tagSet) > 0 && !hasAllTags(v.Tags, tagSet) {
			continue
		}
		results = append(results, Result{Key: v.Key, Value: v.Default, Matched: true})
	}
	return results, nil
}

func hasAllTags(varTags []string, required map[string]bool) bool {
	present := toSet(varTags)
	for tag := range required {
		if !present[tag] {
			return false
		}
	}
	return true
}

func toSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, item := range items {
		if item != "" {
			s[item] = true
		}
	}
	return s
}
