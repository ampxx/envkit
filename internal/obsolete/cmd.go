package obsolete

import (
	"fmt"
	"io"
	"os"

	"github.com/envkit/envkit/internal/config"
)

// RunObsolete loads the config at cfgPath, checks the given target for
// redundant variable declarations, prints a report, and returns a non-nil
// error (or os.ErrProcessDone sentinel) when issues are found.
func RunObsolete(cfgPath, targetName string) error {
	return runObsoleteTo(os.Stdout, cfgPath, targetName)
}

func runObsoleteTo(w io.Writer, cfgPath, targetName string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	report, err := Check(cfg, targetName)
	if err != nil {
		return err
	}

	if !report.HasIssues() {
		fmt.Fprintf(w, "✔  No obsolete variables found in target %q.\n", targetName)
		return nil
	}

	fmt.Fprintf(w, "Obsolete variables in target %q:\n", targetName)
	for _, r := range report.Results {
		fmt.Fprintf(w, "  [%s] %s — %s\n", r.Target, r.Key, r.Reason)
	}
	fmt.Fprintf(w, "\nSummary: %d redundant declaration(s) found.\n", len(report.Results))

	return fmt.Errorf("%d obsolete variable(s) detected", len(report.Results))
}
