package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the comparison between two environment variable sets.
type Result struct {
	OnlyInA  []string
	OnlyInB  []string
	Differing []DiffEntry
	Common   []string
}

// DiffEntry represents a key whose value differs between two targets.
type DiffEntry struct {
	Key    string
	ValueA string
	ValueB string
}

// Compare compares two maps of environment variables and returns a Result.
func Compare(a, b map[string]string) Result {
	result := Result{}

	seen := make(map[string]bool)

	for k, va := range a {
		seen[k] = true
		vb, ok := b[k]
		if !ok {
			result.OnlyInA = append(result.OnlyInA, k)
		} else if va != vb {
			result.Differing = append(result.Differing, DiffEntry{Key: k, ValueA: va, ValueB: vb})
		} else {
			result.Common = append(result.Common, k)
		}
	}

	for k := range b {
		if !seen[k] {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Common)
	sort.Slice(result.Differing, func(i, j int) bool {
		return result.Differing[i].Key < result.Differing[j].Key
	})

	return result
}

// Format returns a human-readable string of the diff result.
func Format(r Result, labelA, labelB string) string {
	var sb strings.Builder

	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(&sb, "Only in %s:\n", labelA)
		for _, k := range r.OnlyInA {
			fmt.Fprintf(&sb, "  - %s\n", k)
		}
	}

	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(&sb, "Only in %s:\n", labelB)
		for _, k := range r.OnlyInB {
			fmt.Fprintf(&sb, "  + %s\n", k)
		}
	}

	if len(r.Differing) > 0 {
		fmt.Fprintf(&sb, "Differing values:\n")
		for _, d := range r.Differing {
			fmt.Fprintf(&sb, "  ~ %s\n    %s: %s\n    %s: %s\n", d.Key, labelA, d.ValueA, labelB, d.ValueB)
		}
	}

	if sb.Len() == 0 {
		return fmt.Sprintf("No differences between %s and %s.\n", labelA, labelB)
	}

	return sb.String()
}
