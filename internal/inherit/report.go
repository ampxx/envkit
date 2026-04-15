package inherit

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable inheritance report to stdout.
func PrintReport(results []Result, src, dst string) {
	printReportTo(os.Stdout, results, src, dst)
}

func printReportTo(w io.Writer, results []Result, src, dst string) {
	fmt.Fprintf(w, "Inheriting from %q → %q\n", src, dst)
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(w, "  SKIP  %-30s (%s)\n", r.Key, r.Reason)
		} else {
			fmt.Fprintf(w, "  OK    %-30s = %s\n", r.Key, r.Value)
		}
	}
	fmt.Fprintln(w, Summary(results))
}

// Summary returns a one-line summary of the inheritance results.
func Summary(results []Result) string {
	inherited, skipped := 0, 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else {
			inherited++
		}
	}
	return fmt.Sprintf("%d inherited, %d skipped", inherited, skipped)
}
