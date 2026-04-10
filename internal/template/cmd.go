package template

import (
	"fmt"
	"os"

	"envkit/internal/config"
)

// RunTemplate is the CLI entry-point for the template command.
//
// Usage: envkit template <config> <target> <template-file> [output-file]
func RunTemplate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envkit template <config> <target> <template-file> [output-file]")
	}

	cfgPath := args[0]
	targetName := args[1]
	tmplPath := args[2]

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	result, err := Generate(cfg, targetName, tmplPath)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if len(args) >= 4 {
		outPath := args[3]
		if err := os.WriteFile(outPath, []byte(result.Output), 0o644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Printf("written to %s (%s)\n", outPath, Summary(result))
		return nil
	}

	Print(result)
	return nil
}
