package trim

import (
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of trimming a single variable.
type Result struct {
	Key      string
	Original string
	Trimmed  string
	Changed  bool
}

// Options controls what trimming is applied.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// TrimQuotes removes surrounding single or double quotes.
	TrimQuotes bool
	// Keys restricts trimming to these keys only (empty = all keys).
	Keys []string
}

// Apply trims variable values in the given target according to opts.
// It returns the set of results and updates cfg in place.
func Apply(cfg *config.Document, targetName string, opts Options) ([]Result, error) {
	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == targetName {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("target %q not found", targetName)
	}

	filter := toSet(opts.Keys)
	var results []Result

	for i := range target.Vars {
		v := &target.Vars[i]
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}
		original := v.Default
		trimmed := original
		if opts.TrimSpace {
			trimmed = strings.TrimSpace(trimmed)
		}
		if opts.TrimQuotes {
			trimmed = trimQuotes(trimmed)
		}
		changed := trimmed != original
		if changed {
			v.Default = trimmed
		}
		results = append(results, Result{
			Key:      v.Key,
			Original: original,
			Trimmed:  trimmed,
			Changed:  changed,
		})
	}
	return results, nil
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
