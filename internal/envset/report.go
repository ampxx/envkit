package envset

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// PrintReport writes a human-readable summary of set results to stdout.
func PrintReport(results []Result) {
	printReportTo(os.Stdout, results)
}

func printReportTo(w io.Writer, results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	for _, r := range results {
		switch {
		case r.Skipped:
			fmt.Fprintf(w, "  SKIP    %s (already set, use --overwrite to replace)\n", r.Key)
		case r.Created:
			fmt.Fprintf(w, "  SET     %s=%s (new)\n", r.Key, r.Value)
		default:
			fmt.Fprintf(w, "  UPDATE  %s=%s\n", r.Key, r.Value)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, Summary(results))
}

// Summary returns a one-line summary string.
func Summary(results []Result) string {
	var set, updated, skipped int
	for _, r := range results {
		switch {
		case r.Skipped:
			skipped++
		case r.Created:
			set++
		default:
			updated++
		}
	}
	return fmt.Sprintf("%d set, %d updated, %d skipped", set, updated, skipped)
}
