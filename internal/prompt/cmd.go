package prompt

import (
	"fmt"
	"os"

	"github.com/user/envkit/internal/config"
	"github.com/user/envkit/internal/validator"
)

// RunFill loads the config for the given target, validates it, and interactively
// prompts the user to fill in any missing required variables, then saves the result.
func RunFill(cfgPath, targetName, envFile string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	target, err := cfg.FindTarget(targetName)
	if err != nil {
		return fmt.Errorf("finding target %q: %w", targetName, err)
	}

	results := validator.Validate(target, envFile)

	var missing []string
	for _, r := range results {
		if r.Status == "missing" {
			missing = append(missing, r.Key)
		}
	}

	if len(missing) == 0 {
		fmt.Println("No missing required variables. Nothing to fill.")
		return nil
	}

	fmt.Printf("Found %d missing required variable(s) for target %q.\n", len(missing), targetName)

	p := New()
	filled, err := p.FillMissing(missing)
	if err != nil {
		return fmt.Errorf("prompting for values: %w", err)
	}

	if len(filled) == 0 {
		fmt.Println("No values provided. Env file unchanged.")
		return nil
	}

	f, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("opening env file %q: %w", envFile, err)
	}
	defer f.Close()

	for k, v := range filled {
		if _, werr := fmt.Fprintf(f, "%s=%s\n", k, v); werr != nil {
			return fmt.Errorf("writing %s to env file: %w", k, werr)
		}
	}

	fmt.Printf("Wrote %d variable(s) to %s\n", len(filled), envFile)
	return nil
}
