package extract

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunExtract is the CLI entry-point for the extract command.
// args: [config-file] --target <name> [--keys k1,k2] [--pattern regex]
func RunExtract(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envkit extract <config> --target <name> [--keys k1,k2] [--pattern regex]")
	}

	cfgPath := args[0]
	opts := Options{}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--target":
			if i+1 >= len(args) {
				return fmt.Errorf("--target requires a value")
			}
			i++
			opts.Target = args[i]
		case "--keys":
			if i+1 >= len(args) {
				return fmt.Errorf("--keys requires a value")
			}
			i++
			for _, k := range strings.Split(args[i], ",") {
				if k = strings.TrimSpace(k); k != "" {
					opts.Keys = append(opts.Keys, k)
				}
			}
		case "--pattern":
			if i+1 >= len(args) {
				return fmt.Errorf("--pattern requires a value")
			}
			i++
			opts.Pattern = args[i]
		}
	}

	if opts.Target == "" {
		return fmt.Errorf("--target is required")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	results, err := Apply(cfg, opts)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no variables matched")
		return nil
	}

	for _, r := range results {
		fmt.Printf("%s=%s\n", r.Key, r.Value)
	}
	fmt.Fprintf(os.Stderr, "extracted %d variable(s) from target %q\n", len(results), opts.Target)
	return nil
}
