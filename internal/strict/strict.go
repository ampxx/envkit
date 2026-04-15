package strict

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of a strict check for a single variable.
type Result struct {
	Target string
	Key    string
	Issue  string
}

// Report is the collection of all strict-mode findings.
type Report struct {
	Results []Result
}

// HasFailures returns true when at least one issue was found.
func (r Report) HasFailures() bool {
	return len(r.Results) > 0
}

// Run performs strict validation on the given target within cfg.
// Strict mode enforces:
//   - no empty values (unless a default is provided)
//   - no keys with lowercase letters
//   - every variable must have a description
func Run(cfg *config.Document, targetName string) (Report, error) {
	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == targetName {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return Report{}, fmt.Errorf("target %q not found", targetName)
	}

	var results []Result
	for _, v := range target.Vars {
		if strings.ToUpper(v.Key) != v.Key {
			results = append(results, Result{
				Target: targetName,
				Key:    v.Key,
				Issue:  "key contains lowercase letters",
			})
		}
		if strings.TrimSpace(v.Description) == "" {
			results = append(results, Result{
				Target: targetName,
				Key:    v.Key,
				Issue:  "missing description",
			})
		}
		if v.Value == "" && v.Default == "" {
			results = append(results, Result{
				Target: targetName,
				Key:    v.Key,
				Issue:  "empty value with no default",
			})
		}
	}
	return Report{Results: results}, nil
}
