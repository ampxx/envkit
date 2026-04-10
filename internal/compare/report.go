package compare

import (
	"fmt"
	"io"
	"os"
)

// PrintReport writes a human-readable diff report to stdout.
func PrintReport(d *TargetDiff) {
	printReportTo(os.Stdout, d)
}

func printReportTo(w io.Writer, d *TargetDiff) {
	fmt.Fprintf(w, "Comparing targets: %s  →  %s\n", d.TargetA, d.TargetB)
	fmt.Fprintf(w, "─────────────────────────────────────────\n")

	if len(d.OnlyInA) > 0 {
		fmt.Fprintf(w, "Only in %s:\n", d.TargetA)
		for _, k := range d.OnlyInA {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(d.OnlyInB) > 0 {
		fmt.Fprintf(w, "Only in %s:\n", d.TargetB)
		for _, k := range d.OnlyInB {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(d.Differing) > 0 {
		fmt.Fprintf(w, "Differing values:\n")
		for _, kd := range d.Differing {
			fmt.Fprintf(w, "  ~ %s\n", kd.Key)
			fmt.Fprintf(w, "      %s: %q\n", d.TargetA, kd.ValueA)
			fmt.Fprintf(w, "      %s: %q\n", d.TargetB, kd.ValueB)
		}
	}

	if len(d.Common) > 0 {
		fmt.Fprintf(w, "Common (identical): %d key(s)\n", len(d.Common))
	}

	fmt.Fprintf(w, "─────────────────────────────────────────\n")
	fmt.Fprintf(w, "%s\n", Summary(d))
}

// Summary returns a one-line summary of the diff.
func Summary(d *TargetDiff) string {
	return fmt.Sprintf(
		"%d only-in-%s  |  %d only-in-%s  |  %d differing  |  %d common",
		len(d.OnlyInA), d.TargetA,
		len(d.OnlyInB), d.TargetB,
		len(d.Differing),
		len(d.Common),
	)
}
