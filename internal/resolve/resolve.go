package resolve

import (
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envkit/internal/config"
)

// Result holds the resolved value for a single variable.
type Result struct {
	Key    string
	Value  string
	Source string // "default", "target", "env", "override"
}

// Options controls resolution behaviour.
type Options struct {
	Target   string
	Override map[string]string // values supplied at call-site (e.g. CLI flags)
	UseEnv   bool              // whether to consult the real process environment
}

// Resolve returns the effective value for every variable defined in cfg for the
// given target, applying the precedence chain:
//
//	override > process env (if UseEnv) > target default > global default
func Resolve(cfg *config.Config, opts Options) ([]Result, error) {
	target := findTarget(cfg, opts.Target)
	if opts.Target != "" && target == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	results := make([]Result, 0, len(cfg.Vars))

	for _, v := range cfg.Vars {
		r := Result{Key: v.Name}

		// 1. global default
		if v.Default != "" {
			r.Value = v.Default
			r.Source = "default"
		}

		// 2. target-level default
		if target != nil {
			for _, tv := range target.Vars {
				if tv.Name == v.Name && tv.Default != "" {
					r.Value = tv.Default
					r.Source = "target"
				}
			}
		}

		// 3. process environment
		if opts.UseEnv {
			if ev, ok := os.LookupEnv(v.Name); ok {
				r.Value = ev
				r.Source = "env"
			}
		}

		// 4. caller override
		if val, ok := opts.Override[v.Name]; ok {
			r.Value = val
			r.Source = "override"
		}

		results = append(results, r)
	}

	return results, nil
}

// ToMap converts a slice of Results into a plain string map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Value
	}
	return m
}

// Format returns a human-readable table of resolved values.
func Format(results []Result, showSource bool) string {
	var sb strings.Builder
	for _, r := range results {
		if showSource {
			fmt.Fprintf(&sb, "%-30s = %-40s  [%s]\n", r.Key, r.Value, r.Source)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", r.Key, r.Value)
		}
	}
	return sb.String()
}

func findTarget(cfg *config.Config, name string) *config.Target {
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == name {
			return &cfg.Targets[i]
		}
	}
	return nil
}
