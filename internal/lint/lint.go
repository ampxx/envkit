package lint

import (
	"fmt"
	"strings"

	"github.com/envkit/envkit/internal/config"
)

// Issue represents a single lint warning or error.
type Issue struct {
	Level   string // "warn" or "error"
	Key     string
	Message string
}

// Result holds all issues found during linting.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue is of level "error".
func (r *Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Level == "error" {
			return true
		}
	}
	return false
}

// Run performs lint checks on the given config and returns a Result.
func Run(cfg *config.Config) *Result {
	result := &Result{}

	seenKeys := map[string]bool{}

	for _, target := range cfg.Targets {
		for _, v := range target.Vars {
			normKey := strings.ToUpper(v.Key)

			// Check for duplicate keys within a target
			qualified := target.Name + ":" + normKey
			if seenKeys[qualified] {
				result.Issues = append(result.Issues, Issue{
					Level:   "error",
					Key:     v.Key,
					Message: fmt.Sprintf("duplicate key %q in target %q", v.Key, target.Name),
				})
			}
			seenKeys[qualified] = true

			// Warn if key is not uppercase
			if v.Key != normKey {
				result.Issues = append(result.Issues, Issue{
					Level:   "warn",
					Key:     v.Key,
					Message: fmt.Sprintf("key %q should be uppercase (expected %q)", v.Key, normKey),
				})
			}

			// Warn if description is missing
			if strings.TrimSpace(v.Description) == "" {
				result.Issues = append(result.Issues, Issue{
					Level:   "warn",
					Key:     v.Key,
					Message: fmt.Sprintf("key %q has no description", v.Key),
				})
			}
		}
	}

	return result
}
