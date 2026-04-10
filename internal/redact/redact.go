package redact

import (
	"regexp"
	"strings"

	"github.com/yourusername/envkit/internal/config"
)

// Rule defines a redaction pattern and its replacement.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// DefaultRules contains common patterns for sensitive values.
var DefaultRules = []Rule{
	{
		Pattern:     regexp.MustCompile(`(?i)(password|secret|token|key|api_key|private)`),
		Replacement: "[REDACTED]",
	},
	{
		Pattern:     regexp.MustCompile(`(?i)(dsn|database_url|connection_string)`),
		Replacement: "[REDACTED]",
	},
}

// Result holds a key and its potentially redacted value.
type Result struct {
	Key      string
	Value    string
	Redacted bool
}

// Apply redacts values in the provided map using the given rules.
// If rules is nil, DefaultRules are used.
func Apply(vars map[string]string, rules []Rule) []Result {
	if rules == nil {
		rules = DefaultRules
	}
	results := make([]Result, 0, len(vars))
	for k, v := range vars {
		redacted := false
		for _, r := range rules {
			if r.Pattern.MatchString(strings.ToLower(k)) {
				v = r.Replacement
				redacted = true
				break
			}
		}
		results = append(results, Result{Key: k, Value: v, Redacted: redacted})
	}
	return results
}

// FromTarget extracts vars for a named target from a config and applies redaction.
func FromTarget(cfg *config.Config, targetName string, rules []Rule) ([]Result, error) {
	for _, t := range cfg.Targets {
		if t.Name == targetName {
			vars := make(map[string]string, len(t.Vars))
			for _, v := range t.Vars {
				vars[v.Key] = v.Default
			}
			return Apply(vars, rules), nil
		}
	}
	return nil, fmt.Errorf("target %q not found", targetName)
}
