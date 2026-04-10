package rotate

import (
	"fmt"
	"os"
	"strings"

	"github.com/envkit/envkit/internal/config"
)

// RunRotate executes the rotate command given CLI arguments.
func RunRotate(cfgPath, target string, keys []string, suffix string, dryRun bool) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	opts := Options{
		Target:  target,
		Keys:    keys,
		Suffix:  suffix,
		DryRun:  dryRun,
	}

	result, err := Rotate(cfg, opts)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	PrintReport(os.Stdout, result)

	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: no changes written)")
		return nil
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// parseSuffixFlag converts a raw suffix flag string into a SuffixFunc.
// Supported forms:
//   - "timestamp" → appends Unix timestamp
//   - "_v2"       → appends literal string
func parseSuffixFlag(raw string) SuffixFunc {
	if strings.ToLower(raw) == "timestamp" {
		return TimestampSuffix()
	}
	return func(_ string) string { return raw }
}
