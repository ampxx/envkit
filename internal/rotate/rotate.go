package rotate

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of a single key rotation.
type Result struct {
	Target  string
	Key     string
	OldVal  string
	NewVal  string
	Rotated bool
	Reason  string
}

// Options controls rotation behaviour.
type Options struct {
	// Keys is the explicit list of keys to rotate; empty means all.
	Keys []string
	// DryRun reports what would change without mutating the config.
	DryRun bool
	// Redact hides old/new values in results (useful for secrets).
	Redact bool
}

// RotateFn is called to produce a new value for a given key + old value.
type RotateFn func(key, oldVal string) (string, error)

// Rotate iterates over vars in every target (or the named targets) and
// calls fn to obtain a replacement value for each matching key.
func Rotate(cfg *config.Config, targetNames []string, opts Options, fn RotateFn) ([]Result, error) {
	if fn == nil {
		return nil, fmt.Errorf("rotate: RotateFn must not be nil")
	}

	wantKey := toSet(opts.Keys)
	wantTarget := toSet(targetNames)

	var results []Result

	for ti := range cfg.Targets {
		t := &cfg.Targets[ti]
		if len(wantTarget) > 0 && !wantTarget[t.Name] {
			continue
		}

		for vi := range t.Vars {
			v := &t.Vars[vi]
			if len(wantKey) > 0 && !wantKey[v.Key] {
				continue
			}

			newVal, err := fn(v.Key, v.Default)
			if err != nil {
				return results, fmt.Errorf("rotate: target %q key %q: %w", t.Name, v.Key, err)
			}

			r := Result{
				Target:  t.Name,
				Key:     v.Key,
				OldVal:  v.Default,
				NewVal:  newVal,
				Rotated: newVal != v.Default,
			}

			if opts.Redact {
				r.OldVal = redactVal(v.Default)
				r.NewVal = redactVal(newVal)
			}

			if r.Rotated && !opts.DryRun {
				v.Default = newVal
			}

			if !r.Rotated {
				r.Reason = "unchanged"
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func redactVal(v string) string {
	if len(v) == 0 {
		return ""
	}
	return strings.Repeat("*", len(v))
}
