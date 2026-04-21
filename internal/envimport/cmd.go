package envimport

import (
	"fmt"
	"strings"

	"envkit/internal/config"

	"github.com/spf13/cobra"
)

// RunImport is the cobra command handler for `envkit import`.
func RunImport(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envkit import <envfile> [flags]")
	}
	envFile := args[0]

	cfgPath, _ := cmd.Flags().GetString("config")
	target, _ := cmd.Flags().GetString("target")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	keysRaw, _ := cmd.Flags().GetString("keys")

	var keys []string
	if keysRaw != "" {
		for _, k := range strings.Split(keysRaw, ",") {
			if trimmed := strings.TrimSpace(k); trimmed != "" {
				keys = append(keys, trimmed)
			}
		}
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	opts := Options{
		Target:    target,
		Overwrite: overwrite,
		DryRun:    dryRun,
		Keys:      keys,
	}

	results, err := Apply(cfg, envFile, opts)
	if err != nil {
		return err
	}

	PrintReport(results, target)

	if dryRun {
		fmt.Println("(dry-run: no changes written)")
		return nil
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}
