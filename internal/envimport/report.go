package envimport

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable summary of import results to stdout.
func PrintReport(results []Result, target string) {
	printReportTo(os.Stdout, results, target)
}

func printReportTo(w io.Writer, results []Result, target string) {
	fmt.Fprintf(w, "Import results for target: %s\n", target)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, r := range results {
		var icon string
		switch r.Status {
		case "imported":
			icon = "+"
		case "updated":
			icon = "~"
		case "skipped":
			icon = "-"
		default:
			icon = "?"
		}
		fmt.Fprintf(w, "  [%s] %-30s %s\n", icon, r.Key, r.Status)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, Summary(results))
}

// Summary returns a one-line summary string.
func Summary(results []Result) string {
	var imported, updated, skipped int
	for _, r := range results {
		switch r.Status {
		case "imported":
			imported++
		case "updated":
			updated++
		case "skipped":
			skipped++
		}
	}
	return fmt.Sprintf("%d imported, %d updated, %d skipped", imported, updated, skipped)
}
