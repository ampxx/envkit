package merge

import (
	"fmt"
	"strings"

	"github.com/user/envkit/internal/config"
)

// Strategy defines how conflicts are resolved during a merge.
type Strategy string

const (
	StrategyOurs   Strategy = "ours"   // keep value from base
	StrategyTheirs Strategy = "theirs" // take value from incoming
	StrategyPrompt Strategy = "prompt" // ask user (handled by caller)
)

// Conflict represents a key whose value differs between two configs.
type Conflict struct {
	Key      string
	BaseVal  string
	OtherVal string
}

// Result holds the merged config and any conflicts encountered.
type Result struct {
	Merged    *config.Config
	Conflicts []Conflict
}

// Merge combines base and other configs using the given strategy.
// Keys only in other are always added. Keys only in base are kept.
// Conflicting keys are resolved per strategy (Prompt conflicts are
// returned unresolved for the caller to handle).
func Merge(base, other *config.Config, strategy Strategy) (*Result, error) {
	if base == nil || other == nil {
		return nil, fmt.Errorf("merge: nil config provided")
	}

	result := &Result{
		Merged: &config.Config{
			Version: base.Version,
			Targets: make(map[string]config.Target),
		},
	}

	// Copy all base targets.
	for name, target := range base.Targets {
		result.Merged.Targets[name] = target
	}

	for name, otherTarget := range other.Targets {
		baseTarget, exists := result.Merged.Targets[name]
		if !exists {
			result.Merged.Targets[name] = otherTarget
			continue
		}

		merged, conflicts := mergeVars(baseTarget.Vars, otherTarget.Vars, strategy)
		result.Conflicts = append(result.Conflicts, conflicts...)
		baseTarget.Vars = merged
		result.Merged.Targets[name] = baseTarget
	}

	return result, nil
}

func mergeVars(base, other []config.VarDef, strategy Strategy) ([]config.VarDef, []Conflict) {
	index := make(map[string]int, len(base))
	result := make([]config.VarDef, len(base))
	copy(result, base)
	for i, v := range result {
		index[v.Name] = i
	}

	var conflicts []Conflict
	for _, ov := range other {
		idx, exists := index[ov.Name]
		if !exists {
			result = append(result, ov)
			continue
		}
		if result[idx].Default == ov.Default {
			continue
		}
		switch strategy {
		case StrategyTheirs:
			result[idx].Default = ov.Default
		case StrategyOurs:
			// keep base value — no-op
		case StrategyPrompt:
			conflicts = append(conflicts, Conflict{
				Key:      ov.Name,
				BaseVal:  result[idx].Default,
				OtherVal: ov.Default,
			})
		}
	}
	return result, conflicts
}

// ApplyResolutions updates the merged config with user-resolved conflicts.
func ApplyResolutions(result *Result, resolutions map[string]string) {
	for target, t := range result.Merged.Targets {
		for i, v := range t.Vars {
			if val, ok := resolutions[v.Name]; ok {
				t.Vars[i].Default = strings.TrimSpace(val)
			}
		}
		result.Merged.Targets[target] = t
	}
}
