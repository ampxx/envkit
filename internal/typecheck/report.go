package typecheck

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes the typecheck report to stdout.
func PrintReport(r Report) {
	printReportTo(os.Stdout, r)
}

func printReportTo(w io.Writer, r Report) {
	fmt.Fprintf(w, "Type check — target: %s\n", r.Target)
	for _, res := range r.Results {
		switch res.Status {
		case StatusOK:
			fmt.Fprintf(w, "  ✓ %-30s [%s]\n", res.Key, res.Expected)
		case StatusFail:
			fmt.Fprintf(w, "  ✗ %-30s [%s] %s\n", res.Key, res.Expected, res.Reason)
		case StatusSkipped:
			if res.Reason != "" {
				fmt.Fprintf(w, "  - %-30s [%s] (%s)\n", res.Key, res.Expected, res.Reason)
			} else {
				fmt.Fprintf(w, "  - %-30s [no type]\n", res.Key)
			}
		}
	}
	fmt.Fprintln(w, Summary(r))
}

// Summary returns a one-line summary of the report.
func Summary(r Report) string {
	ok, fail, skip := 0, 0, 0
	for _, res := range r.Results {
		switch res.Status {
		case StatusOK:
			ok++
		case StatusFail:
			fail++
		case StatusSkipped:
			skip++
		}
	}
	return fmt.Sprintf("%d ok, %d failed, %d skipped", ok, fail, skip)
}
