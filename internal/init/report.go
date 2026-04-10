package init

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PrintReport writes a human-readable summary of the scaffold operation.
func PrintReport(opts Options, tmpl Template) {
	printReportTo(os.Stdout, opts, tmpl)
}

func printReportTo(w io.Writer, opts Options, tmpl Template) {
	fmt.Fprintf(w, "✔ Scaffolded envkit.yaml in %s\n", opts.Dir)
	fmt.Fprintf(w, "  Template : %s — %s\n", tmpl.Name, tmpl.Description)
	fmt.Fprintf(w, "  Targets  : %s\n", strings.Join(tmpl.Targets, ", "))
	fmt.Fprintf(w, "  Variables:\n")
	for _, v := range tmpl.Vars {
		req := ""
		if v.Required {
			req = " (required)"
		}
		fmt.Fprintf(w, "    - %s%s\n", v.Key, req)
	}
	fmt.Fprintf(w, "\nNext steps:\n")
	fmt.Fprintf(w, "  1. Edit envkit.yaml and fill in your variable definitions.\n")
	fmt.Fprintf(w, "  2. Run `envkit validate` to check your config.\n")
	fmt.Fprintf(w, "  3. Run `envkit doctor` to diagnose your environment files.\n")
}

// Summary returns a one-line description of the scaffold result.
func Summary(opts Options) string {
	return fmt.Sprintf("created envkit.yaml in %s using template %q", opts.Dir, opts.Template)
}
