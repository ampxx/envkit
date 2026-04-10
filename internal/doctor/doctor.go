package doctor

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
	"envkit/internal/validator"
)

// CheckResult represents the outcome of a single health check.
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
}

// Report holds all check results for a doctor run.
type Report struct {
	Target  string
	Results []CheckResult
}

// HasFailures returns true if any check failed.
func (r *Report) HasFailures() bool {
	for _, res := range r.Results {
		if !res.Passed {
			return true
		}
	}
	return false
}

// Run performs all health checks for the given target against the provided env file.
func Run(cfg *config.Config, targetName string, envFile string) (*Report, error) {
	target, err := config.FindTarget(cfg, targetName)
	if err != nil {
		return nil, fmt.Errorf("doctor: %w", err)
	}

	report := &Report{Target: targetName}

	// Check 1: env file exists
	report.Results = append(report.Results, checkFileExists(envFile))

	// Check 2: no empty keys in config
	report.Results = append(report.Results, checkNoEmptyKeys(cfg))

	// Check 3: validate env against target rules
	if _, err := os.Stat(envFile); err == nil {
		results, err := validator.Validate(target, envFile)
		if err != nil {
			return nil, fmt.Errorf("doctor: validation error: %w", err)
		}
		report.Results = append(report.Results, checkValidation(results))
	}

	// Check 4: no duplicate variable names in target
	report.Results = append(report.Results, checkNoDuplicateVars(target))

	return report, nil
}

func checkFileExists(path string) CheckResult {
	_, err := os.Stat(path)
	if err != nil {
		return CheckResult{Name: "env file exists", Passed: false, Message: fmt.Sprintf("%s not found", path)}
	}
	return CheckResult{Name: "env file exists", Passed: true, Message: path + " found"}
}

func checkNoEmptyKeys(cfg *config.Config) CheckResult {
	for _, t := range cfg.Targets {
		for _, v := range t.Vars {
			if strings.TrimSpace(v.Key) == "" {
				return CheckResult{Name: "no empty keys", Passed: false, Message: "empty key found in target " + t.Name}
			}
		}
	}
	return CheckResult{Name: "no empty keys", Passed: true, Message: "all keys are non-empty"}
}

func checkValidation(results []validator.Result) CheckResult {
	if validator.HasFailures(results) {
		return CheckResult{Name: "env validation", Passed: false, Message: fmt.Sprintf("%d rule(s) failed", countFailed(results))}
	}
	return CheckResult{Name: "env validation", Passed: true, Message: "all rules passed"}
}

func countFailed(results []validator.Result) int {
	n := 0
	for _, r := range results {
		if !r.Passed {
			n++
		}
	}
	return n
}

func checkNoDuplicateVars(target *config.Target) CheckResult {
	seen := map[string]bool{}
	for _, v := range target.Vars {
		if seen[v.Key] {
			return CheckResult{Name: "no duplicate vars", Passed: false, Message: "duplicate key: " + v.Key}
		}
		seen[v.Key] = true
	}
	return CheckResult{Name: "no duplicate vars", Passed: true, Message: "no duplicates found"}
}
