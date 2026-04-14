package summary

import (
	"fmt"
	"sort"

	"envkit/internal/config"
)

// Result holds the summary information for a single target.
type Result struct {
	Target      string
	TotalVars   int
	Required    int
	Optional    int
	Sensitive   int
	Tagged      int
	WithDefault int
}

// Report is the full summary across all (or selected) targets.
type Report struct {
	Results []Result
}

// Run produces a summary Report for the given config.
// If targets is non-empty only those targets are included;
// passing nil/empty summarises every target.
func Run(cfg *config.Document, targets []string) (Report, error) {
	wantSet := toSet(targets)

	var results []Result
	for _, t := range cfg.Targets {
		if len(wantSet) > 0 {
			if _, ok := wantSet[t.Name]; !ok {
				continue
			}
		}

		r := Result{Target: t.Name, TotalVars: len(t.Vars)}
		for _, v := range t.Vars {
			if v.Required {
				r.Required++
			} else {
				r.Optional++
			}
			if v.Sensitive {
				r.Sensitive++
			}
			if len(v.Tags) > 0 {
				r.Tagged++
			}
			if v.Default != "" {
				r.WithDefault++
			}
		}
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Target < results[j].Target
	})

	if len(results) == 0 && len(wantSet) > 0 {
		return Report{}, fmt.Errorf("no matching targets found")
	}

	return Report{Results: results}, nil
}

func toSet(items []string) map[string]struct{} {
	m := make(map[string]struct{}, len(items))
	for _, s := range items {
		m[s] = struct{}{}
	}
	return m
}
