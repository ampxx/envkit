package validator

import (
	"fmt"
	"io"
	"strings"
)

const (
	passIcon = "✓"
	failIcon = "✗"
)

// PrintReport writes a human-readable validation report to w.
func PrintReport(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No variables validated.")
		return
	}

	passed := 0
	failed := 0

	fmt.Fprintln(w, strings.Repeat("-", 50))
	for _, r := range results {
		icon := passIcon
		if !r.Passed {
			icon = failIcon
			failed++
		} else {
			passed++
		}
		fmt.Fprintf(w, "  %s  %-30s %s\n", icon, r.Key, r.Message)
	}
	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintf(w, "  Passed: %d  Failed: %d\n", passed, failed)
}

// Summary returns a short one-line summary string.
func Summary(results []Result) string {
	total := len(results)
	failed := 0
	for _, r := range results {
		if !r.Passed {
			failed++
		}
	}
	if failed == 0 {
		return fmt.Sprintf("All %d variable(s) passed validation.", total)
	}
	return fmt.Sprintf("%d of %d variable(s) failed validation.", failed, total)
}

// FailedResults returns only the results that did not pass validation.
// This is useful when callers want to act on or display only failures.
func FailedResults(results []Result) []Result {
	var failed []Result
	for _, r := range results {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	return failed
}
