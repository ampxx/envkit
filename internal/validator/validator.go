package validator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Required bool
	Pattern  string
	Allowed  []string
}

// Result holds the outcome of a single variable validation.
type Result struct {
	Key     string
	Passed  bool
	Message string
}

// Validate checks environment variables against the provided rules map.
// It reads values from the supplied env map (or os.Environ if nil).
func Validate(rules map[string]Rule, env map[string]string) []Result {
	if env == nil {
		env = envToMap()
	}

	var results []Result

	for key, rule := range rules {
		val, exists := env[key]

		if rule.Required && !exists {
			results = append(results, Result{
				Key:     key,
				Passed:  false,
				Message: fmt.Sprintf("%s is required but not set", key),
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, val)
			if err != nil || !matched {
				results = append(results, Result{
					Key:     key,
					Passed:  false,
					Message: fmt.Sprintf("%s value %q does not match pattern %q", key, val, rule.Pattern),
				})
				continue
			}
		}

		if len(rule.Allowed) > 0 && !contains(rule.Allowed, val) {
			results = append(results, Result{
				Key:     key,
				Passed:  false,
				Message: fmt.Sprintf("%s value %q is not in allowed set [%s]", key, val, strings.Join(rule.Allowed, ", ")),
			})
			continue
		}

		results = append(results, Result{Key: key, Passed: true, Message: "ok"})
	}

	return results
}

// HasFailures returns true if any result failed.
func HasFailures(results []Result) bool {
	for _, r := range results {
		if !r.Passed {
			return true
		}
	}
	return false
}

func envToMap() map[string]string {
	m := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
