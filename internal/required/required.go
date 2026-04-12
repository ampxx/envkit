package required

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds the check outcome for a single variable.
type Result struct {
	Key     string
	Target  string
	Present bool
	HasDefault bool
}

// Report is the full output of a required-vars check.
type Report struct {
	Results []Result
}

// HasFailures returns true when any required variable is absent and has no default.
func (r *Report) HasFailures() bool {
	for _, res := range r.Results {
		if !res.Present && !res.HasDefault {
			return true
		}
	}
	return false
}

// Check inspects the supplied env map against the config for targetName and
// returns a Report listing every required variable and whether it is satisfied.
func Check(cfg *config.Document, targetName string, env map[string]string) (*Report, error) {
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

	report := &Report{}
	for _, v := range target.Vars {
		if !v.Required {
			continue
		}
		_, present := env[v.Key]
		report.Results = append(report.Results, Result{
			Key:        v.Key,
			Target:     targetName,
			Present:    present,
			HasDefault: v.Default != "",
		})
	}
	return report, nil
}
