package merge

import (
	"fmt"
	"io"
	"strings"
)

// PrintConflicts writes a human-readable summary of merge conflicts to w.
func PrintConflicts(w io.Writer, conflicts []Conflict) {
	if len(conflicts) == 0 {
		fmt.Fprintln(w, "No conflicts.")
		return
	}
	fmt.Fprintf(w, "%d conflict(s) found:\n", len(conflicts))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range conflicts {
		fmt.Fprintf(w, "  KEY   : %s\n", c.Key)
		fmt.Fprintf(w, "  BASE  : %s\n", c.BaseVal)
		fmt.Fprintf(w, "  OTHER : %s\n", c.OtherVal)
		fmt.Fprintln(w, strings.Repeat("-", 40))
	}
}

// SummaryLine returns a one-line merge summary suitable for audit logs.
func SummaryLine(result *Result) string {
	if result == nil {
		return "merge: no result"
	}
	targets := make([]string, 0, len(result.Merged.Targets))
	for k := range result.Merged.Targets {
		targets = append(targets, k)
	}
	return fmt.Sprintf("merged %d target(s) [%s], %d conflict(s) resolved",
		len(targets), strings.Join(targets, ", "), len(result.Conflicts))
}
