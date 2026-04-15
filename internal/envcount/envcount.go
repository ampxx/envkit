package envcount

import (
	"fmt"
	"sort"

	"envkit/internal/config"
)

// Result holds the variable count for a single target.
type Result struct {
	Target   string
	Total    int
	Required int
	Optional int
	Sensitive int
}

// Report is the full output of a Count run.
type Report struct {
	Results []Result
}

// Count returns a Report describing how many variables each target declares.
// If targets is non-empty only those targets are included; otherwise all
// targets in the config are counted.
func Count(cfg *config.Document, targets []string) (Report, error) {
	wantSet := toSet(targets)

	var results []Result
	for _, t := range cfg.Targets {
		if len(wantSet) > 0 {
			if _, ok := wantSet[t.Name]; !ok {
				continue
			}
		}

		r := Result{Target: t.Name, Total: len(t.Vars)}
		for _, v := range t.Vars {
			if v.Required {
				r.Required++
			} else {
				r.Optional++
			}
			if v.Sensitive {
				r.Sensitive++
			}
		}
		results = append(results, r)
	}

	if len(wantSet) > 0 && len(results) == 0 {
		return Report{}, fmt.Errorf("no matching targets found")
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Target < results[j].Target
	})

	return Report{Results: results}, nil
}

// GrandTotal returns the sum of all variable counts across all results.
func (r Report) GrandTotal() int {
	total := 0
	for _, res := range r.Results {
		total += res.Total
	}
	return total
}

func toSet(ss []string) map[string]struct{} {
	m := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		m[s] = struct{}{}
	}
	return m
}
