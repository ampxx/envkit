package render

import (
	"fmt"
	"os"

	"envkit/internal/config"
)

// RunRender is the CLI entry-point for the `envkit render` command.
//
// Usage:
//
//	envkit render <target> --template <file> [--out <file>]
func RunRender(configPath, targetName, templatePath, outPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var result *Result
	if templatePath != "" {
		result, err = RenderFile(cfg, targetName, templatePath)
	} else {
		return fmt.Errorf("--template is required")
	}
	if err != nil {
		return err
	}

	if outPath != "" {
		if err := os.WriteFile(outPath, []byte(result.Output), 0o644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Fprintf(os.Stdout, "rendered %q → %s\n", targetName, outPath)
		return nil
	}

	fmt.Print(result.Output)
	return nil
}
