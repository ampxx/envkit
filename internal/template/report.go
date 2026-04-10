package template

import (
	"fmt"
	"io"
	"os"
)

// Print writes the template result to stdout with a header.
func Print(r *Result) {
	printTo(os.Stdout, r)
}

// PrintTo writes the template result to an arbitrary writer.
func printTo(w io.Writer, r *Result) {
	fmt.Fprintf(w, "# envkit template — target: %s\n", r.Target)
	fmt.Fprintln(w, r.Output)
}

// Summary returns a one-line description of the result.
func Summary(r *Result) string {
	lines := 0
	for _, c := range r.Output {
		if c == '\n' {
			lines++
		}
	}
	return fmt.Sprintf("target=%s lines=%d", r.Target, lines)
}
