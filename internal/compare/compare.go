package compare

import (
	"fmt"
	"sort"

	"github.com/user/envkit/internal/config"
)

// TargetDiff holds the result of comparing two deployment targets.
type TargetDiff struct {
	TargetA   string
	TargetB   string
	OnlyInA   []string
	OnlyInB   []string
	Differing []KeyDiff
	Common    []string
}

// KeyDiff represents a key whose value differs between two targets.
type KeyDiff struct {
	Key    string
	ValueA string
	ValueB string
}

// Targets compares the variable definitions of two named targets within a config.
func Targets(cfg *config.Config, targetA, targetB string) (*TargetDiff, error) {
	aVars, err := varsForTarget(cfg, targetA)
	if err != nil {
		return nil, err
	}
	bVars, err := varsForTarget(cfg, targetB)
	if err != nil {
		return nil, err
	}

	result := &TargetDiff{TargetA: targetA, TargetB: targetB}

	aKeys := keySet(aVars)
	bKeys := keySet(bVars)

	for k := range aKeys {
		if _, ok := bKeys[k]; !ok {
			result.OnlyInA = append(result.OnlyInA, k)
		} else {
			av := aVars[k]
			bv := bVars[k]
			if av != bv {
				result.Differing = append(result.Differing, KeyDiff{Key: k, ValueA: av, ValueB: bv})
			} else {
				result.Common = append(result.Common, k)
			}
		}
	}
	for k := range bKeys {
		if _, ok := aKeys[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Common)
	sort.Slice(result.Differing, func(i, j int) bool {
		return result.Differing[i].Key < result.Differing[j].Key
	})

	return result, nil
}

// HasDifferences returns true if the diff contains any keys that are exclusive
// to one target or have differing values between the two targets.
func (d *TargetDiff) HasDifferences() bool {
	return len(d.OnlyInA) > 0 || len(d.OnlyInB) > 0 || len(d.Differing) > 0
}

func varsForTarget(cfg *config.Config, name string) (map[string]string, error) {
	for _, t := range cfg.Targets {
		if t.Name == name {
			m := make(map[string]string, len(t.Vars))
			for _, v := range t.Vars {
				m[v.Key] = v.Default
			}
			return m, nil
		}
	}
	return nil, fmt.Errorf("target %q not found", name)
}

func keySet(m map[string]string) map[string]struct{} {
	s := make(map[string]struct{}, len(m))
	for k := range m {
		s[k] = struct{}{}
	}
	return s
}
