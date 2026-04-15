package envdiff

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PrintReport writes a human-readable diff report to stdout.
func PrintReport(r *Report, showIdentical bool) {
	printReportTo(os.Stdout, r, showIdentical)
}

func printReportTo(w io.Writer, r *Report, showIdentical bool) {
	fmt.Fprintf(w, "envdiff: %s  ↔  %s\n", r.TargetA, r.TargetB)
	fmt.Fprintln(w, strings.Repeat("─", 50))

	for _, res := range r.Results {
		switch res.Kind {
		case KindOnlyInA:
			fmt.Fprintf(w, "  ← only in %-12s  %s=%s\n", r.TargetA, res.Key, res.ValueA)
		case KindOnlyInB:
			fmt.Fprintf(w, "  → only in %-12s  %s=%s\n", r.TargetB, res.Key, res.ValueB)
		case KindDiffering:
			fmt.Fprintf(w, "  ~ %-20s  [%s] %s  ≠  [%s] %s\n",
				res.Key, r.TargetA, res.ValueA, r.TargetB, res.ValueB)
		case KindIdentical:
			if showIdentical {
				fmt.Fprintf(w, "  = %-20s  %s\n", res.Key, res.ValueA)
			}
		}
	}

	fmt.Fprintln(w, strings.Repeat("─", 50))
	fmt.Fprintln(w, Summary(r))
}

// Summary returns a one-line summary string for the report.
func Summary(r *Report) string {
	var onlyA, onlyB, diff, same int
	for _, res := range r.Results {
		switch res.Kind {
		case KindOnlyInA:
			onlyA++
		case KindOnlyInB:
			onlyB++
		case KindDiffering:
			diff++
		case KindIdentical:
			same++
		}
	}
	return fmt.Sprintf("only-in-%s: %d  only-in-%s: %d  differing: %d  identical: %d",
		r.TargetA, onlyA, r.TargetB, onlyB, diff, same)
}
