package promote

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable promotion report to stdout.
func PrintReport(results []Result, from, to string) {
	printReportTo(os.Stdout, results, from, to)
}

func printReportTo(w io.Writer, results []Result, from, to string) {
	fmt.Fprintf(w, "Promoting variables: %s → %s\n\n", from, to)
	promoted := 0
	skipped := 0
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(w, "  SKIP  %-30s  %s\n", r.Key, r.Reason)
			skipped++
		} else if r.OldValue != "" {
			fmt.Fprintf(w, "  UPDATE %-30s  %q → %q\n", r.Key, r.OldValue, r.NewValue)
			promoted++
		} else {
			fmt.Fprintf(w, "  ADD   %-30s  %q\n", r.Key, r.NewValue)
			promoted++
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Summary: %d promoted, %d skipped\n", promoted, skipped)
}

// Summary returns a one-line summary string.
func Summary(results []Result) string {
	promoted, skipped := 0, 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else {
			promoted++
		}
	}
	return fmt.Sprintf("%d promoted, %d skipped", promoted, skipped)
}
