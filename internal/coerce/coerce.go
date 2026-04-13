package coerce

import (
	"fmt"
	"strconv"
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of coercing a single variable.
type Result struct {
	Key      string
	Original string
	Coerced  string
	Type     string
	Changed  bool
	Skipped  bool
}

// Options controls coercion behaviour.
type Options struct {
	Target string
	Keys   []string // empty means all keys
	DryRun bool
}

// Apply coerces variable values in the given target to their declared types.
func Apply(cfg *config.Document, env map[string]string, opts Options) ([]Result, error) {
	target := findTarget(cfg, opts.Target)
	if target == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	filter := toSet(opts.Keys)
	var results []Result

	for _, v := range target.Vars {
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}
		raw, ok := env[v.Key]
		if !ok {
			results = append(results, Result{Key: v.Key, Skipped: true})
			continue
		}
		coerced, err := coerceValue(raw, v.Type)
		if err != nil {
			results = append(results, Result{Key: v.Key, Original: raw, Skipped: true})
			continue
		}
		changed := coerced != raw
		if changed && !opts.DryRun {
			env[v.Key] = coerced
		}
		results = append(results, Result{
			Key:      v.Key,
			Original: raw,
			Coerced:  coerced,
			Type:     v.Type,
			Changed:  changed,
		})
	}
	return results, nil
}

func coerceValue(val, typ string) (string, error) {
	switch strings.ToLower(typ) {
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return "", fmt.Errorf("cannot coerce %q to bool", val)
		}
		return strconv.FormatBool(b), nil
	case "int":
		i, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot coerce %q to int", val)
		}
		return strconv.FormatInt(i, 10), nil
	case "float":
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return "", fmt.Errorf("cannot coerce %q to float", val)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	default:
		return val, nil
	}
}

func findTarget(cfg *config.Document, name string) *config.Target {
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == name {
			return &cfg.Targets[i]
		}
	}
	return nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
