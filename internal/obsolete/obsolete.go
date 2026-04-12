package obsolete

import (
	"fmt"

	"github.com/envkit/envkit/internal/config"
)

// Result holds information about a single obsolete variable.
type Result struct {
	Target string
	Key    string
	Reason string
}

// Report is the full output of an obsolete check.
type Report struct {
	Results []Result
}

// HasIssues returns true when at least one obsolete variable was found.
func (r Report) HasIssues() bool {
	return len(r.Results) > 0
}

// Check inspects the given config for variables that are declared in a target
// but are already present (with an identical value) in the global defaults,
// making the target-level declaration redundant.
func Check(cfg *config.Document, targetName string) (Report, error) {
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

	// Build a map of global defaults.
	globals := make(map[string]string)
	for _, v := range cfg.Vars {
		globals[v.Key] = v.Default
	}

	var results []Result
	for _, v := range target.Vars {
		globalDefault, exists := globals[v.Key]
		if !exists {
			continue
		}
		if v.Default == globalDefault {
			results = append(results, Result{
				Target: targetName,
				Key:    v.Key,
				Reason: fmt.Sprintf("value %q duplicates global default", v.Default),
			})
		}
	}

	return Report{Results: results}, nil
}
