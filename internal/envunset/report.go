package envunset

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable report of unset results to stdout.
func PrintReport(results []Result) {
	printReportTo(os.Stdout, results)
}

func printReportTo(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No keys targeted for removal.")
		return
	}
	for _, r := range results {
		if r.Found {
			fmt.Fprintf(w, "  [removed] %s  (target: %s)\n", r.Key, r.Target)
		} else {
			fmt.Fprintf(w, "  [not found] %s  (target: %s)\n", r.Key, r.Target)
		}
	}
	fmt.Fprintln(w, Summary(results))
}

// Summary returns a one-line summary string for the given results.
func Summary(results []Result) string {
	removed, missing := 0, 0
	for _, r := range results {
		if r.Found {
			removed++
		} else {
			missing++
		}
	}
	return fmt.Sprintf("%d removed, %d not found", removed, missing)
}
