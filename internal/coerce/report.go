package coerce

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable coercion report to stdout.
func PrintReport(results []Result) {
	printReportTo(os.Stdout, results)
}

func printReportTo(w io.Writer, results []Result) {
	for _, r := range results {
		switch {
		case r.Skipped:
			fmt.Fprintf(w, "  SKIP  %s\n", r.Key)
		case r.Changed:
			fmt.Fprintf(w, "  CHANGED  %-24s  %q -> %q  (%s)\n", r.Key, r.Original, r.Coerced, r.Type)
		default:
			fmt.Fprintf(w, "  OK     %-24s  %q\n", r.Key, r.Coerced)
		}
	}
}

// Summary returns a one-line summary of coercion results.
func Summary(results []Result) string {
	changed, skipped, ok := 0, 0, 0
	for _, r := range results {
		switch {
		case r.Skipped:
			skipped++
		case r.Changed:
			changed++
		default:
			ok++
		}
	}
	return fmt.Sprintf("%d changed, %d ok, %d skipped", changed, ok, skipped)
}
