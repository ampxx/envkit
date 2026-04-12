package unused

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds information about a variable that is declared but not used
// in any real environment file across the specified targets.
type Result struct {
	Target string
	Key    string
	Reason string
}

// Report is the full output of a Check run.
type Report struct {
	Results []Result
}

// Check scans the config for variables that have no default value and are not
// referenced in any of the provided env maps (keyed by target name).
func Check(cfg *config.Document, envMaps map[string]map[string]string) Report {
	var results []Result

	for _, target := range cfg.Targets {
		envMap := envMaps[target.Name]

		for _, v := range target.Vars {
			_, inEnv := envMap[v.Key]
			hasDefault := v.Default != ""

			if !inEnv && !hasDefault {
				results = append(results, Result{
					Target: target.Name,
					Key:    v.Key,
					Reason: "not present in env file and has no default",
				})
			}
		}

		// Also flag keys present in env but not declared in config.
		for k := range envMap {
			if !declaredInTarget(target, k) {
				results = append(results, Result{
					Target: target.Name,
					Key:    k,
					Reason: "present in env file but not declared in config",
				})
			}
		}
	}

	return Report{Results: results}
}

// HasIssues returns true when the report contains at least one result.
func (r Report) HasIssues() bool {
	return len(r.Results) > 0
}

// Summary returns a short human-readable summary line.
func (r Report) Summary() string {
	if !r.HasIssues() {
		return "no unused or undeclared variables found"
	}
	return fmt.Sprintf("%d issue(s) found", len(r.Results))
}

func declaredInTarget(target config.Target, key string) bool {
	for _, v := range target.Vars {
		if v.Key == key {
			return true
		}
	}
	return false
}
