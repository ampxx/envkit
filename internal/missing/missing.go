// Package missing identifies environment variables that are defined in a
// config target but absent from a live .env file.
package missing

import (
	"fmt"

	"envkit/internal/config"
	"envkit/internal/env"
)

// Result holds information about a single missing variable.
type Result struct {
	Key         string
	Required    bool
	Description string
}

// Report is the output of a missing-variable check.
type Report struct {
	Target  string
	Results []Result
}

// HasMissing returns true when at least one variable is absent.
func (r Report) HasMissing() bool { return len(r.Results) > 0 }

// RequiredMissing returns only the results that are marked required.
func (r Report) RequiredMissing() []Result {
	out := make([]Result, 0)
	for _, res := range r.Results {
		if res.Required {
			out = append(out, res)
		}
	}
	return out
}

// Check compares the variables declared in targetName against the key/value
// pairs loaded from envFile and returns a Report of anything absent.
func Check(cfg *config.Document, targetName, envFile string) (Report, error) {
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

	parsed, err := env.ParseFile(envFile)
	if err != nil {
		return Report{}, fmt.Errorf("parsing env file: %w", err)
	}
	live := env.ToMap(parsed)

	report := Report{Target: targetName}
	for _, v := range target.Vars {
		if _, ok := live[v.Key]; !ok {
			report.Results = append(report.Results, Result{
				Key:         v.Key,
				Required:    v.Required,
				Description: v.Description,
			})
		}
	}
	return report, nil
}
