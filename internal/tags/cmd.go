package tags

import (
	"fmt"
	"strings"

	"envkit/internal/config"
)

// RunTags is the CLI entry-point for the `tags` command.
//
// Usage:
//
//	envkit tags --target=<name> [--add=t1,t2] [--remove=t3] [--keys=K1,K2]
//	envkit tags list --target=<name>
func RunTags(cfgPath, target, addCSV, removeCSV, keysCSV string, listOnly bool) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if listOnly {
		return runList(cfg, target)
	}
	return runApply(cfg, cfgPath, target, addCSV, removeCSV, keysCSV)
}

func runList(cfg *config.Document, target string) error {
	tags, err := List(cfg, target)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		fmt.Printf("No tags found for target %q\n", target)
		return nil
	}
	fmt.Printf("Tags for target %q:\n", target)
	for _, t := range tags {
		fmt.Printf("  - %s\n", t)
	}
	return nil
}

func runApply(cfg *config.Document, cfgPath, target, addCSV, removeCSV, keysCSV string) error {
	opts := Options{
		Target: target,
		Keys:   parseCSV(keysCSV),
		Add:    parseCSV(addCSV),
		Remove: parseCSV(removeCSV),
	}

	results, err := Apply(cfg, opts)
	if err != nil {
		return err
	}

	updated := 0
	for _, r := range results {
		if r.Action == "updated" {
			updated++
			fmt.Printf("  updated  %s  [%s]\n", r.Key, strings.Join(r.Tags, ", "))
		}
	}

	if updated == 0 {
		fmt.Println("No changes.")
		return nil
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	fmt.Printf("\n%d variable(s) updated in target %q.\n", updated, target)
	return nil
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
