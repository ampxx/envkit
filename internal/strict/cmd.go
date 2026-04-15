package strict

import (
	"fmt"
	"io"
	"os"

	"envkit/internal/config"
)

// RunStrict loads the config from cfgPath, runs strict checks against
// targetName and prints a human-readable report to stdout.
// It returns a non-nil error when the config cannot be loaded or the
// target is not found. When issues are found it prints them and returns
// an error so the caller (CLI) can exit with a non-zero status.
func RunStrict(cfgPath, targetName string) error {
	return runStrictTo(os.Stdout, cfgPath, targetName)
}

func runStrictTo(w io.Writer, cfgPath, targetName string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	report, err := Run(cfg, targetName)
	if err != nil {
		return err
	}

	if !report.HasFailures() {
		fmt.Fprintf(w, "strict: target %q passed all checks\n", targetName)
		return nil
	}

	fmt.Fprintf(w, "strict: %d issue(s) found in target %q\n\n", len(report.Results), targetName)
	for _, r := range report.Results {
		fmt.Fprintf(w, "  [%s] %s\n", r.Key, r.Issue)
	}
	fmt.Fprintln(w)
	return fmt.Errorf("strict checks failed with %d issue(s)", len(report.Results))
}
