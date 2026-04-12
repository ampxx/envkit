package copy

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable copy report to stdout.
func PrintReport(results []Result, from, to string) {
	printReportTo(os.Stdout, results, from, to)
}

func printReportTo(w io.Writer, results []Result, from, to string) {
	fmt.Fprintf(w, "Copy: %s → %s\n", from, to)
	if len(results) == 0 {
		fmt.Fprintln(w, "  no variables matched")
		return
	}
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(w, "  ~ %-30s skipped (%s)\n", r.Key, r.Reason)
		} else {
			fmt.Fprintf(w, "  + %-30s copied\n", r.Key)
		}
	}
	fmt.Fprintln(w, Summary(results))
}

// Summary returns a one-line summary of copy results.
func Summary(results []Result) string {
	copied, skipped := 0, 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else {
			copied++
		}
	}
	return fmt.Sprintf("%d copied, %d skipped", copied, skipped)
}
