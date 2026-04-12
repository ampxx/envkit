package typecheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"envkit/internal/config"
)

// ResultStatus represents the outcome of a type check.
type ResultStatus string

const (
	StatusOK      ResultStatus = "ok"
	StatusFail    ResultStatus = "fail"
	StatusSkipped ResultStatus = "skipped"
)

// Result holds the outcome of checking a single variable.
type Result struct {
	Key      string
	Expected string
	Actual   string
	Status   ResultStatus
	Reason   string
}

// Report is the full output of a typecheck run.
type Report struct {
	Target  string
	Results []Result
}

var emailRe = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
var urlRe = regexp.MustCompile(`^https?://`)

// Run checks the actual values in env against the declared types in cfg for the given target.
func Run(cfg *config.Document, target string, env map[string]string) (Report, error) {
	var tgt *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == target {
			tgt = &cfg.Targets[i]
			break
		}
	}
	if tgt == nil {
		return Report{}, fmt.Errorf("target %q not found", target)
	}

	report := Report{Target: target}
	for _, v := range tgt.Vars {
		if v.Type == "" {
			report.Results = append(report.Results, Result{
				Key: v.Key, Expected: "any", Actual: env[v.Key], Status: StatusSkipped,
			})
			continue
		}
		actual, ok := env[v.Key]
		if !ok {
			report.Results = append(report.Results, Result{
				Key: v.Key, Expected: v.Type, Actual: "", Status: StatusSkipped, Reason: "not set",
			})
			continue
		}
		status, reason := checkType(v.Type, actual)
		report.Results = append(report.Results, Result{
			Key: v.Key, Expected: v.Type, Actual: actual, Status: status, Reason: reason,
		})
	}
	return report, nil
}

func checkType(typ, val string) (ResultStatus, string) {
	switch strings.ToLower(typ) {
	case "int", "integer":
		if _, err := strconv.Atoi(val); err != nil {
			return StatusFail, fmt.Sprintf("%q is not an integer", val)
		}
	case "float", "number":
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return StatusFail, fmt.Sprintf("%q is not a number", val)
		}
	case "bool", "boolean":
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return StatusFail, fmt.Sprintf("%q is not a boolean", val)
		}
	case "url":
		if !urlRe.MatchString(val) {
			return StatusFail, fmt.Sprintf("%q is not a valid URL", val)
		}
	case "email":
		if !emailRe.MatchString(val) {
			return StatusFail, fmt.Sprintf("%q is not a valid email", val)
		}
	}
	return StatusOK, ""
}

// HasFailures returns true if any result in the report has StatusFail.
func HasFailures(r Report) bool {
	for _, res := range r.Results {
		if res.Status == StatusFail {
			return true
		}
	}
	return false
}
