package dedupe

import (
	"fmt"

	"envkit/internal/config"
)

// Result holds information about a single deduplication action.
type Result struct {
	Target string
	Key    string
	Kept   string
	Removed []string
}

// Report is the full output of a dedupe run.
type Report struct {
	Results []Result
	Total   int
}

// Apply scans the given target for duplicate variable keys and removes all but
// the last occurrence (last-write-wins). If targetName is empty, all targets
// are scanned.
func Apply(cfg *config.Document, targetName string) (Report, error) {
	var report Report

	for i, target := range cfg.Targets {
		if targetName != "" && target.Name != targetName {
			continue
		}

		seen := map[string][]int{} // key -> slice of indices
		for idx, v := range target.Vars {
			seen[v.Key] = append(seen[v.Key], idx)
		}

		// Collect indices to remove (all but the last occurrence).
		removeSet := map[int]bool{}
		for key, indices := range seen {
			if len(indices) <= 1 {
				continue
			}
			kept := indices[len(indices)-1]
			var removedVals []string
			for _, idx := range indices[:len(indices)-1] {
				removeSet[idx] = true
				removedVals = append(removedVals, target.Vars[idx].Value)
			}
			report.Results = append(report.Results, Result{
				Target:  target.Name,
				Key:     key,
				Kept:    target.Vars[kept].Value,
				Removed: removedVals,
			})
			report.Total += len(removedVals)
		}

		if len(removeSet) == 0 {
			continue
		}

		filtered := cfg.Targets[i].Vars[:0]
		for idx, v := range target.Vars {
			if !removeSet[idx] {
				filtered = append(filtered, v)
			}
		}
		cfg.Targets[i].Vars = filtered
	}

	if targetName != "" {
		found := false
		for _, t := range cfg.Targets {
			if t.Name == targetName {
				found = true
				break
			}
		}
		if !found {
			return Report{}, fmt.Errorf("target %q not found", targetName)
		}
	}

	return report, nil
}
