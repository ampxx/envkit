package copy

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunCopy implements the `envkit copy <from> <to>` command.
func RunCopy(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envkit copy <from-target> <to-target> [--keys=A,B] [--overwrite] [--file=envkit.yaml]")
	}

	from := args[0]
	to := args[1]

	filePath := flags["file"]
	if filePath == "" {
		filePath = "envkit.yaml"
	}

	cfg, err := config.Load(filePath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	var keys []string
	if raw := flags["keys"]; raw != "" {
		for _, k := range strings.Split(raw, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	_, overwrite := flags["overwrite"]

	opts := Options{
		Keys:      keys,
		Overwrite: overwrite,
	}

	results, updated, err := Copy(cfg, from, to, opts)
	if err != nil {
		return err
	}

	PrintReport(results, from, to)

	if err := config.Save(updated, filePath); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Fprintf(os.Stdout, "saved to %s\n", filePath)
	return nil
}
