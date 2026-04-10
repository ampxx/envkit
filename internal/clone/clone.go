package clone

import (
	"fmt"

	"github.com/envkit/envkit/internal/config"
)

// Result holds the outcome of a clone operation.
type Result struct {
	SourceTarget string
	DestTarget   string
	Copied       int
	Skipped      int
	Overwritten  int
}

// Options controls clone behaviour.
type Options struct {
	Overwrite bool
	Keys      []string // if non-empty, only clone these keys
}

// Clone duplicates all variable definitions from srcName into dstName within
// the provided config. If dstName does not exist it is created.
func Clone(cfg *config.Config, srcName, dstName string, opts Options) (Result, error) {
	src := findTarget(cfg, srcName)
	if src == nil {
		return Result{}, fmt.Errorf("source target %q not found", srcName)
	}

	dst := findTarget(cfg, dstName)
	if dst == nil {
		cfg.Targets = append(cfg.Targets, config.Target{Name: dstName})
		dst = &cfg.Targets[len(cfg.Targets)-1]
	}

	filter := toSet(opts.Keys)

	res := Result{SourceTarget: srcName, DestTarget: dstName}

	for _, v := range src.Vars {
		if len(filter) > 0 && !filter[v.Key] {
			res.Skipped++
			continue
		}

		existing := findVar(dst, v.Key)
		if existing != nil {
			if !opts.Overwrite {
				res.Skipped++
				continue
			}
			*existing = v
			res.Overwritten++
			continue
		}

		dst.Vars = append(dst.Vars, v)
		res.Copied++
	}

	return res, nil
}

func findTarget(cfg *config.Config, name string) *config.Target {
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == name {
			return &cfg.Targets[i]
		}
	}
	return nil
}

func findVar(t *config.Target, key string) *config.VarDef {
	for i := range t.Vars {
		if t.Vars[i].Key == key {
			return &t.Vars[i]
		}
	}
	return nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
