package scope

import (
	"fmt"
	"os"
	"strings"

	"envkit/internal/config"
)

// RunScope is the CLI entry point for the scope command.
func RunScope(cfgPath, target, prefix string, keys, exclude []string, quiet bool) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	opts := Options{
		Target:  target,
		Keys:    keys,
		Exclude: exclude,
	}

	res, err := Apply(cfg, opts)
	if err != nil {
		return err
	}

	if prefix != "" {
		res = Prefix(res, prefix)
	}

	if quiet {
		for k, v := range res.Vars {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	PrintReport(res)
	return nil
}

// parseCSV splits a comma-separated string into a trimmed slice.
func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
