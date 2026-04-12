package patch

import (
	"fmt"

	"envkit/internal/config"
)

// Op represents a patch operation type.
type Op string

const (
	OpSet    Op = "set"
	OpUnset  Op = "unset"
	OpRename Op = "rename"
)

// Change describes a single patch operation.
type Change struct {
	Op      Op
	Key     string
	Value   string // used by set
	NewKey  string // used by rename
}

// Result records the outcome of applying a Change.
type Result struct {
	Change  Change
	Applied bool
	Note    string
}

// Apply applies a slice of Changes to the named target within cfg.
// It returns the list of results and mutates cfg in place.
func Apply(cfg *config.Config, target string, changes []Change) ([]Result, error) {
	var tgt *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == target {
			tgt = &cfg.Targets[i]
			break
		}
	}
	if tgt == nil {
		return nil, fmt.Errorf("target %q not found", target)
	}

	results := make([]Result, 0, len(changes))
	for _, ch := range changes {
		r := applyOne(tgt, ch)
		results = append(results, r)
	}
	return results, nil
}

func applyOne(tgt *config.Target, ch Change) Result {
	switch ch.Op {
	case OpSet:
		for i := range tgt.Vars {
			if tgt.Vars[i].Key == ch.Key {
				old := tgt.Vars[i].Default
				tgt.Vars[i].Default = ch.Value
				return Result{Change: ch, Applied: true, Note: fmt.Sprintf("updated %q: %q -> %q", ch.Key, old, ch.Value)}
			}
		}
		tgt.Vars = append(tgt.Vars, config.VarDef{Key: ch.Key, Default: ch.Value})
		return Result{Change: ch, Applied: true, Note: fmt.Sprintf("added %q=%q", ch.Key, ch.Value)}

	case OpUnset:
		for i := range tgt.Vars {
			if tgt.Vars[i].Key == ch.Key {
				tgt.Vars = append(tgt.Vars[:i], tgt.Vars[i+1:]...)
				return Result{Change: ch, Applied: true, Note: fmt.Sprintf("removed %q", ch.Key)}
			}
		}
		return Result{Change: ch, Applied: false, Note: fmt.Sprintf("key %q not found", ch.Key)}

	case OpRename:
		for i := range tgt.Vars {
			if tgt.Vars[i].Key == ch.Key {
				tgt.Vars[i].Key = ch.NewKey
				return Result{Change: ch, Applied: true, Note: fmt.Sprintf("renamed %q -> %q", ch.Key, ch.NewKey)}
			}
		}
		return Result{Change: ch, Applied: false, Note: fmt.Sprintf("key %q not found", ch.Key)}

	default:
		return Result{Change: ch, Applied: false, Note: fmt.Sprintf("unknown op %q", ch.Op)}
	}
}
