package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the diff outcome for a single key.
type Result struct {
	Key     string
	ValueA  string
	ValueB  string
	OnlyInA bool
	OnlyInB bool
}

// Compare returns a map of keys that differ between two env maps.
func Compare(a, b map[string]string) map[string]Result {
	out := make(map[string]Result)

	for k, va := range a {
		if vb, ok := b[k]; !ok {
			out[k] = Result{Key: k, ValueA: va, OnlyInA: true}
		} else if va != vb {
			out[k] = Result{Key: k, ValueA: va, ValueB: vb}
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			out[k] = Result{Key: k, ValueB: vb, OnlyInB: true}
		}
	}

	return out
}

// Format returns a formatted string representation of the diff results.
func Format(result map[string]Result) string {
	if len(result) == 0 {
		return "(no differences)"
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		r := result[k]
		switch {
		case r.OnlyInA:
			fmt.Fprintf(&sb, "- %s=%q\n", k, r.ValueA)
		case r.OnlyInB:
			fmt.Fprintf(&sb, "+ %s=%q\n", k, r.ValueB)
		default:
			fmt.Fprintf(&sb, "~ %s: %q -> %q\n", k, r.ValueA, r.ValueB)
		}
	}
	return sb.String()
}
