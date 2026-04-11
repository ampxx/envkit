package mask

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"envkit/internal/config"
)

// RunMask loads the config, applies masking for the given target, and prints
// the results to stdout. extraRules is a slice of "KEY=strategy" strings.
func RunMask(cfgPath, target string, extraRules []string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	extras, err := parseRules(extraRules)
	if err != nil {
		return err
	}

	results, err := Apply(cfg, target, extras)
	if err != nil {
		return err
	}

	printMaskResults(results)
	return nil
}

func parseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY=strategy", r)
		}
		strategy := strings.ToLower(strings.TrimSpace(parts[1]))
		switch strategy {
		case "full", "partial", "hash", "none":
		default:
			return nil, fmt.Errorf("unknown strategy %q for key %q", strategy, parts[0])
		}
		rules = append(rules, Rule{Key: strings.TrimSpace(parts[0]), Strategy: strategy})
	}
	return rules, nil
}

func printMaskResults(results []Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tSTRATEGY\tVALUE")
	fmt.Fprintln(w, "---\t--------\t-----")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Key, r.Strategy, r.Masked)
	}
	w.Flush()
}
