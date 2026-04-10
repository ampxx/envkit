package promote

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunPromote is the CLI entry point for the promote command.
// Args: [from-target] [to-target] [flags]
func RunPromote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envkit promote <from> <to> [--overwrite] [--keys KEY1,KEY2]")
	}

	from := args[0]
	to := args[1]

	opts := Options{}
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--overwrite":
			opts.Overwrite = true
		case "--keys":
			if i+1 >= len(args) {
				return fmt.Errorf("--keys requires a comma-separated list of key names")
			}
			i++
			for _, k := range strings.Split(args[i], ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					opts.Keys = append(opts.Keys, k)
				}
			}
		default:
			return fmt.Errorf("unknown flag: %s", args[i])
		}
	}

	cfg, err := config.Load("envkit.yaml")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	results, err := Promote(cfg, from, to, opts)
	if err != nil {
		return err
	}

	PrintReport(results, from, to)

	if err := config.Save(cfg, "envkit.yaml"); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Config saved.")
	return nil
}
