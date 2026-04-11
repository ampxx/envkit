package mask

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// Rule defines how a key should be masked.
type Rule struct {
	Key      string
	Strategy string // "full", "partial", "hash"
}

// Result holds the masked output for a single variable.
type Result struct {
	Key      string
	Original string
	Masked   string
	Strategy string
}

// Apply masks environment variable values for the given target according to
// the variable definitions (sensitive flag) and any extra rules provided.
func Apply(cfg *config.Document, target string, extras []Rule) ([]Result, error) {
	var t *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == target {
			t = &cfg.Targets[i]
			break
		}
	}
	if t == nil {
		return nil, fmt.Errorf("target %q not found", target)
	}

	extraMap := make(map[string]string, len(extras))
	for _, r := range extras {
		extraMap[strings.ToUpper(r.Key)] = r.Strategy
	}

	var results []Result
	for _, v := range t.Vars {
		key := strings.ToUpper(v.Key)
		strategy := "none"
		if v.Sensitive {
			strategy = "full"
		}
		if s, ok := extraMap[key]; ok {
			strategy = s
		}
		masked := maskValue(v.Default, strategy)
		results = append(results, Result{
			Key:      v.Key,
			Original: v.Default,
			Masked:   masked,
			Strategy: strategy,
		})
	}
	return results, nil
}

func maskValue(val, strategy string) string {
	switch strategy {
	case "full":
		return "********"
	case "partial":
		if len(val) <= 4 {
			return "****"
		}
		return val[:2] + strings.Repeat("*", len(val)-4) + val[len(val)-2:]
	case "hash":
		return fmt.Sprintf("[sha:%x]", simpleHash(val))
	default:
		return val
	}
}

func simpleHash(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}
