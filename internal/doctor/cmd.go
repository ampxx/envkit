package doctor

import (
	"fmt"
	"os"

	"envkit/internal/config"
)

// RunDoctor loads the config and runs all health checks, printing results.
func RunDoctor(cfgPath, targetName, envFile string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("doctor: failed to load config: %w", err)
	}

	report, err := Run(cfg, targetName, envFile)
	if err != nil {
		return err
	}

	PrintReport(report)
	fmt.Println(Summary(report))

	if report.HasFailures() {
		os.Exit(1)
	}
	return nil
}
