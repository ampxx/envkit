package inherit

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunInherit is the CLI entry-point for the inherit command.
// Usage: envkit inherit <config> <src-target> <dst-target> [--keys=A,B] [--overwrite] [--save]
func RunInherit(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: inherit <config> <src-target> <dst-target> [--keys=...] [--overwrite] [--save]")
	}
	cfgPath := args[0]
	srcTarget := args[1]
	dstTarget := args[2]

	var keys []string
	overwrite := false
	save := false

	for _, arg := range args[3:] {
		switch {
		case strings.HasPrefix(arg, "--keys="):
			raw := strings.TrimPrefix(arg, "--keys=")
			for _, k := range strings.Split(raw, ",") {
				if k = strings.TrimSpace(k); k != "" {
					keys = append(keys, k)
				}
			}
		case arg == "--overwrite":
			overwrite = true
		case arg == "--save":
			save = true
		}
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	results, err := Apply(cfg, srcTarget, dstTarget, Options{
		Keys:      keys,
		Overwrite: overwrite,
	})
	if err != nil {
		return err
	}

	PrintReport(results, srcTarget, dstTarget)

	if save {
		if err := config.Save(cfg, cfgPath); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Config saved to %s\n", cfgPath)
	}
	return nil
}
