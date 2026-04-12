package flatten

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunFlatten is the CLI entry-point for the flatten command.
func RunFlatten(cfgPath, target, prefix, separator string, keys []string, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	results, err := Apply(cfg, Options{
		Target:    target,
		Prefix:    prefix,
		Separator: separator,
		Keys:      keys,
	})
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(out, "# no variables found")
		return nil
	}

	header := "global"
	if target != "" {
		header = target
	}
	fmt.Fprintf(out, "# flattened vars — target: %s\n", header)

	for _, r := range results {
		fmt.Fprintf(out, "%s=%s\n", r.Key, quoteIfNeeded(r.Value))
	}

	fmt.Fprintf(out, "# total: %d\n", len(results))
	return nil
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n") {
		return `"` + v + `"`
	}
	return v
}
