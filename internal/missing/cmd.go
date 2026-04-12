package missing

import (
	"fmt"
	"io"
	"os"

	"envkit/internal/config"
)

// RunMissing loads the config at cfgPath, checks targetName against envFile,
// prints a human-readable report, and returns a non-nil error when required
// variables are absent (so the caller can exit non-zero).
func RunMissing(cfgPath, targetName, envFile string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	rep, err := Check(cfg, targetName, envFile)
	if err != nil {
		return err
	}

	printReport(os.Stdout, rep)

	if len(rep.RequiredMissing()) > 0 {
		return fmt.Errorf("%d required variable(s) missing from %s", len(rep.RequiredMissing()), envFile)
	}
	return nil
}

func printReport(w io.Writer, rep Report) {
	if !rep.HasMissing() {
		fmt.Fprintf(w, "✓ All variables present for target %q\n", rep.Target)
		return
	}
	fmt.Fprintf(w, "Missing variables for target %q:\n", rep.Target)
	for _, r := range rep.Results {
		requiredTag := ""
		if r.Required {
			requiredTag = " [required]"
		}
		desc := ""
		if r.Description != "" {
			desc = "  — " + r.Description
		}
		fmt.Fprintf(w, "  - %s%s%s\n", r.Key, requiredTag, desc)
	}
	fmt.Fprintf(w, "\nSummary: %d missing (%d required)\n",
		len(rep.Results), len(rep.RequiredMissing()))
}
