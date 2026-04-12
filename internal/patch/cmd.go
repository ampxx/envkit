package patch

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// RunPatch loads the config at cfgPath, applies patch operations expressed as
// key=value strings (set), key- (unset), or oldKey>newKey (rename) against
// the given target, saves the result and prints a summary.
func RunPatch(cfgPath, target string, ops []string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	changes, err := parseOps(ops)
	if err != nil {
		return err
	}

	results, err := Apply(cfg, target, changes)
	if err != nil {
		return err
	}

	applied, skipped := 0, 0
	for _, r := range results {
		if r.Applied {
			applied++
			fmt.Printf("  ✔ %s\n", r.Note)
		} else {
			skipped++
			fmt.Printf("  ✗ %s\n", r.Note)
		}
	}
	fmt.Printf("\npatch: %d applied, %d skipped\n", applied, skipped)

	if applied > 0 {
		if err := config.Save(cfg, cfgPath); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
	}
	return nil
}

// parseOps converts raw string tokens into Change structs.
// Supported syntax:
//   KEY=VALUE  -> set
//   KEY-       -> unset
//   OLD>NEW    -> rename
func parseOps(ops []string) ([]Change, error) {
	changes := make([]Change, 0, len(ops))
	for _, op := range ops {
		switch {
		case strings.Contains(op, "="):
			parts := strings.SplitN(op, "=", 2)
			changes = append(changes, Change{Op: OpSet, Key: parts[0], Value: parts[1]})
		case strings.HasSuffix(op, "-"):
			changes = append(changes, Change{Op: OpUnset, Key: strings.TrimSuffix(op, "-")})
		case strings.Contains(op, ">"):
			parts := strings.SplitN(op, ">", 2)
			changes = append(changes, Change{Op: OpRename, Key: parts[0], NewKey: parts[1]})
		default:
			return nil, fmt.Errorf("unrecognised patch op %q (use KEY=VAL, KEY-, or OLD>NEW)", op)
		}
	}
	return changes, nil
}
