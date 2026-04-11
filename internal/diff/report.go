package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// PrintReport writes a human-readable diff report to stdout.
func PrintReport(result map[string]Result) {
	printReportTo(os.Stdout, result)
}

func printReportTo(w io.Writer, result map[string]Result) {
	if len(result) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		r := result[k]
		switch {
		case r.OnlyInA:
			fmt.Fprintf(w, "  - %-30s [only in A] %q\n", k, r.ValueA)
		case r.OnlyInB:
			fmt.Fprintf(w, "  + %-30s [only in B] %q\n", k, r.ValueB)
		default:
			fmt.Fprintf(w, "  ~ %-30s A=%q  B=%q\n", k, r.ValueA, r.ValueB)
		}
	}
}

// Summary returns a one-line summary of the diff result.
func Summary(result map[string]Result) string {
	onlyA, onlyB, differ := 0, 0, 0
	for _, r := range result {
		switch {
		case r.OnlyInA:
			onlyA++
		case r.OnlyInB:
			onlyB++
		default:
			differ++
		}
	}
	if onlyA == 0 && onlyB == 0 && differ == 0 {
		return "files are identical"
	}
	return fmt.Sprintf("%d only-in-A, %d only-in-B, %d differing", onlyA, onlyB, differ)
}
