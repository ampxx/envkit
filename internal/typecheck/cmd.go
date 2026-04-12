package typecheck

import (
	"fmt"
	"os"

	"envkit/internal/config"
	"envkit/internal/env"
)

// RunTypeCheck loads the config and env file, runs type checking, and prints the report.
func RunTypeCheck(cfgPath, envFile, target string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var envMap map[string]string
	if envFile != "" {
		entries, err := env.ParseFile(envFile)
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}
		envMap = env.ToMap(entries)
	} else {
		envMap = captureOSEnv()
	}

	report, err := Run(cfg, target, envMap)
	if err != nil {
		return err
	}

	PrintReport(report)

	if HasFailures(report) {
		os.Exit(1)
	}
	return nil
}

// captureOSEnv reads the current process environment into a map.
func captureOSEnv() map[string]string {
	m := make(map[string]string)
	for _, pair := range os.Environ() {
		for i := 0; i < len(pair); i++ {
			if pair[i] == '=' {
				m[pair[:i]] = pair[i+1:]
				break
			}
		}
	}
	return m
}
