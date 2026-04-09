package audit

import (
	"fmt"
	"os"
	"text/tabwriter"
)

const defaultLogPath = ".envkit-audit.log"

// RunAuditLog prints the audit log entries to stdout.
func RunAuditLog(logPath string) error {
	if logPath == "" {
		logPath = defaultLogPath
	}

	entries, err := ReadAll(logPath)
	if err != nil {
		return fmt.Errorf("failed to read audit log: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No audit log entries found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tEVENT\tTARGET\tUSER\tSUCCESS\tDETAILS")
	fmt.Fprintln(w, "---------\t-----\t------\t----\t-------\t-------")

	for _, e := range entries {
		status := "ok"
		if !e.Success {
			status = "FAIL"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Event,
			e.Target,
			e.User,
			status,
			e.Details,
		)
	}

	return w.Flush()
}

// RecordEvent is a convenience helper for other commands to log an audit event.
func RecordEvent(event EventType, target string, success bool, details string) {
	l := New(defaultLogPath)
	if err := l.Log(event, target, success, details); err != nil {
		// audit logging is best-effort; do not fail the main command
		fmt.Fprintf(os.Stderr, "warning: audit log failed: %v\n", err)
	}
}
