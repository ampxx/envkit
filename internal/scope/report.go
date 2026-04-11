package scope

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// PrintReport writes the scoped result to stdout.
func PrintReport(r Result) {
	printReportTo(os.Stdout, r)
}

func printReportTo(w io.Writer, r Result) {
	fmt.Fprintf(w, "Scope: target=%s  vars=%d\n", r.Target, len(r.Vars))
	fmt.Fprintln(w, "")

	keys := make([]string, 0, len(r.Vars))
	for k := range r.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "  %-30s = %s\n", k, r.Vars[k])
	}
}

// Summary returns a one-line description of the result.
func Summary(r Result) string {
	return fmt.Sprintf("target=%s vars=%d", r.Target, len(r.Vars))
}
