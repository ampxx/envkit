package intersect

import "github.com/your-org/envkit/internal/config"

// Result holds the intersection analysis between two targets.
type Result struct {
	Target string
	OtherTarget string
	CommonKeys []string
	OnlyInA    []string
	OnlyInB    []string
}

// Apply computes the intersection and differences of variable keys
// declared in targetA and targetB within the given config.
func Apply(cfg *config.Document, targetA, targetB string) (Result, error) {
	keysA, err := keysForTarget(cfg, targetA)
	if err != nil {
		return Result{}, err
	}
	keysB, err := keysForTarget(cfg, targetB)
	if err != nil {
		return Result{}, err
	}

	setA := toSet(keysA)
	setB := toSet(keysB)

	var common, onlyA, onlyB []string

	for k := range setA {
		if setB[k] {
			common = append(common, k)
		} else {
			onlyA = append(onlyA, k)
		}
	}
	for k := range setB {
		if !setA[k] {
			onlyB = append(onlyB, k)
		}
	}

	sortStrings(common)
	sortStrings(onlyA)
	sortStrings(onlyB)

	return Result{
		Target:      targetA,
		OtherTarget: targetB,
		CommonKeys:  common,
		OnlyInA:    onlyA,
		OnlyInB:    onlyB,
	}, nil
}

func keysForTarget(cfg *config.Document, name string) ([]string, error) {
	for _, t := range cfg.Targets {
		if t.Name == name {
			keys := make([]string, 0, len(t.Vars))
			for _, v := range t.Vars {
				keys = append(keys, v.Key)
			}
			return keys, nil
		}
	}
	return nil, fmt.Errorf("target %q not found", name)
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func sortStrings(s []string) {
	sort.Strings(s)
}
