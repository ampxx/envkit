package required

import (
	"fmt"
	"os"

	"envkit/internal/config"
	"envkit/internal/env"
)

// RunRequired checks required variables for a target against an env file.
// configPath is the envkit config file; target is the deployment target name;
// envFile is the .env file to read actual values from.
func RunRequired(configPath, target, envFile string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	entries, err := env.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}
	environment := env.ToMap(entries)

	report, err := Check(cfg, target, environment)
	if err != nil {
		return err
	}

	printReport(report)

	if report.HasFailures() {
		return fmt.Errorf("required variable check failed")
	}
	return nil
}

func printReport(r *Report) {
	for _, res := range r.Results {
		status := "OK"
		if !res.Present {
			if res.HasDefault {
				status = "MISSING (has default)"
			} else {
				status = "MISSING"
			}
		}
		fmt.Fprintf(os.Stdout, "  %-30s %s\n", res.Key, status)
	}

	total := len(r.Results)
	failed := 0
	for _, res := range r.Results {
		if !res.Present && !res.HasDefault {
			failed++
		}
	}
	fmt.Fprintf(os.Stdout, "\n%d required variable(s) checked, %d missing\n", total, failed)
}
