package merge

import (
	"fmt"
	"os"

	"github.com/user/envkit/internal/audit"
	"github.com/user/envkit/internal/config"
	"github.com/user/envkit/internal/prompt"
)

// RunMerge loads two config files, merges them, and writes the result.
// baseFile is treated as the authoritative source; otherFile is merged in.
func RunMerge(baseFile, otherFile, outFile, strategyStr string, auditLog *audit.Logger) error {
	base, err := config.Load(baseFile)
	if err != nil {
		return fmt.Errorf("loading base config %q: %w", baseFile, err)
	}

	other, err := config.Load(otherFile)
	if err != nil {
		return fmt.Errorf("loading other config %q: %w", otherFile, err)
	}

	strategy := Strategy(strategyStr)
	switch strategy {
	case StrategyOurs, StrategyTheirs, StrategyPrompt:
		// valid
	default:
		return fmt.Errorf("unknown merge strategy %q (choose: ours, theirs, prompt)", strategyStr)
	}

	result, err := Merge(base, other, strategy)
	if err != nil {
		return err
	}

	if len(result.Conflicts) > 0 {
		resolutions, err := resolveConflicts(result.Conflicts)
		if err != nil {
			return err
		}
		ApplyResolutions(result, resolutions)
	}

	if err := config.Save(result.Merged, outFile); err != nil {
		return fmt.Errorf("saving merged config: %w", err)
	}

	if auditLog != nil {
		_ = auditLog.Log("merge", fmt.Sprintf("%s + %s -> %s [%s]", baseFile, otherFile, outFile, strategy))
	}

	fmt.Fprintf(os.Stdout, "✔ merged config written to %s\n", outFile)
	return nil
}

func resolveConflicts(conflicts []Conflict) (map[string]string, error) {
	p := prompt.New(os.Stdin, os.Stdout)
	resolutions := make(map[string]string, len(conflicts))

	fmt.Fprintln(os.Stdout, "\nConflicts detected — please resolve each:")
	for _, c := range conflicts {
		fmt.Fprintf(os.Stdout, "  Key: %s\n  Base:  %s\n  Other: %s\n", c.Key, c.BaseVal, c.OtherVal)
		useOurs, err := p.AskConfirm(fmt.Sprintf("Keep base value %q?", c.BaseVal))
		if err != nil {
			return nil, err
		}
		if useOurs {
			resolutions[c.Key] = c.BaseVal
		} else {
			resolutions[c.Key] = c.OtherVal
		}
	}
	return resolutions, nil
}
