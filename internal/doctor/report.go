package doctor

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	passIcon = "✓"
	failIcon = "✗"
)

// PrintReport writes a formatted doctor report to stdout.
func PrintReport(report *Report) {
	printReportTo(os.Stdout, report)
}

func printReportTo(w io.Writer, report *Report) {
	fmt.Fprintf(w, "Doctor report for target: %s\n", report.Target)
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, r := range report.Results {
		icon := passIcon
		if !r.Passed {
			icon = failIcon
		}
		fmt.Fprintf(w, "  %s  %-25s %s\n", icon, r.Name, r.Message)
	}
	fmt.Fprintln(w)
	if report.HasFailures() {
		fmt.Fprintln(w, "Status: UNHEALTHY — some checks failed.")
	} else {
		fmt.Fprintln(w, "Status: HEALTHY — all checks passed.")
	}
}

// Summary returns a one-line summary string.
func Summary(report *Report) string {
	total := len(report.Results)
	passed := 0
	for _, r := range report.Results {
		if r.Passed {
			passed++
		}
	}
	status := "HEALTHY"
	if report.HasFailures() {
		status = "UNHEALTHY"
	}
	return fmt.Sprintf("[%s] %s: %d/%d checks passed", status, report.Target, passed, total)
}
