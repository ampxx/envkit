package normalize

import (
	"strings"

	"envkit/internal/config"
)

// Rule defines how a variable value should be normalized.
type Rule struct {
	Uppercase bool
	Lowercase bool
	TrimSpace bool
	TrimQuotes bool
}

// Result holds the outcome of normalizing a single variable.
type Result struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
}

// Report holds all normalization results for a target.
type Report struct {
	Target  string
	Results []Result
}

// Apply normalizes variable values in the given target according to the provided rule.
// If keys is non-empty, only those keys are normalized.
func Apply(cfg *config.Document, targetName string, rule Rule, keys []string) (Report, error) {
	target, err := findTarget(cfg, targetName)
	if err != nil {
		return Report{}, err
	}

	keySet := toSet(keys)
	report := Report{Target: targetName}

	for i, v := range target.Vars {
		if len(keySet) > 0 && !keySet[v.Key] {
			continue
		}

		old := v.Default
		newVal := old

		if rule.TrimSpace {
			newVal = strings.TrimSpace(newVal)
		}
		if rule.TrimQuotes {
			newVal = strings.Trim(newVal, `"'`)
		}
		if rule.Uppercase {
			newVal = strings.ToUpper(newVal)
		} else if rule.Lowercase {
			newVal = strings.ToLower(newVal)
		}

		changed := newVal != old
		if changed {
			target.Vars[i].Default = newVal
		}

		report.Results = append(report.Results, Result{
			Key:      v.Key,
			OldValue: old,
			NewValue: newVal,
			Changed:  changed,
		})
	}

	return report, nil
}

func findTarget(cfg *config.Document, name string) (*config.Target, error) {
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == name {
			return &cfg.Targets[i], nil
		}
	}
	return nil, fmt.Errorf("target %q not found", name)
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
