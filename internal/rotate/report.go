package rotate

import (
	"fmt"
	"io"
)

// PrintReport writes a human-readable rotation report to w.
func PrintReport(w io.Writer, r Result) {
	if len(r.Rotated) == 0 && len(r.Skipped) == 0 {
		fmt.Fprintln(w, "no variables matched rotation criteria")
		return
	}

	if len(r.Rotated) > 0 {
		fmt.Fprintln(w, "rotated:")
		for _, rv := range r.Rotated {
			fmt.Fprintf(w, "  %-30s  %s → %s\n", rv.Key, redactVal(rv.OldValue), redactVal(rv.NewValue))
		}
	}

	if len(r.Skipped) > 0 {
		fmt.Fprintln(w, "skipped:")
		for _, key := range r.Skipped {
			fmt.Fprintf(w, "  %s\n", key)
		}
	}

	fmt.Fprintln(w, Summary(r))
}

// Summary returns a one-line summary string for the rotation result.
func Summary(r Result) string {
	return fmt.Sprintf("summary: %d rotated, %d skipped", len(r.Rotated), len(r.Skipped))
}
