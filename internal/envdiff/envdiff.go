package envdiff

import (
	"sort"

	"github.com/your-org/envkit/internal/config"
)

// ResultKind describes how a variable differs between two targets.
type ResultKind string

const (
	KindOnlyInA  ResultKind = "only_in_a"
	KindOnlyInB  ResultKind = "only_in_b"
	KindDiffering ResultKind = "differing"
	KindIdentical ResultKind = "identical"
)

// Result holds the comparison outcome for a single key.
type Result struct {
	Key    string
	Kind   ResultKind
	ValueA string
	ValueB string
}

// Report is the full set of results for a target-pair comparison.
type Report struct {
	TargetA string
	TargetB string
	Results []Result
}

// Compare compares the variables of two named targets within a config document.
// Values from environment variable defaults are used when a target does not
// override a key.
func Compare(doc *config.Document, targetA, targetB string) (*Report, error) {
	mapA, err := varsForTarget(doc, targetA)
	if err != nil {
		return nil, err
	}
	mapB, err := varsForTarget(doc, targetB)
	if err != nil {
		return nil, err
	}

	keys := unionKeys(mapA, mapB)
	sort.Strings(keys)

	report := &Report{TargetA: targetA, TargetB: targetB}
	for _, k := range keys {
		va, inA := mapA[k]
		vb, inB := mapB[k]
		switch {
		case inA && !inB:
			report.Results = append(report.Results, Result{Key: k, Kind: KindOnlyInA, ValueA: va})
		case !inA && inB:
			report.Results = append(report.Results, Result{Key: k, Kind: KindOnlyInB, ValueB: vb})
		case va != vb:
			report.Results = append(report.Results, Result{Key: k, Kind: KindDiffering, ValueA: va, ValueB: vb})
		default:
			report.Results = append(report.Results, Result{Key: k, Kind: KindIdentical, ValueA: va, ValueB: vb})
		}
	}
	return report, nil
}

func varsForTarget(doc *config.Document, name string) (map[string]string, error) {
	m := make(map[string]string)
	for _, v := range doc.Vars {
		if v.Default != "" {
			m[v.Key] = v.Default
		}
	}
	for i := range doc.Targets {
		if doc.Targets[i].Name == name {
			for _, v := range doc.Targets[i].Vars {
				m[v.Key] = v.Default
			}
			return m, nil
		}
	}
	return nil, fmt.Errorf("envdiff: target %q not found", name)
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
