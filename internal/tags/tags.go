package tags

import (
	"fmt"
	"sort"
	"strings"

	"envkit/internal/config"
)

// Result holds the outcome of a tag operation for a single variable.
type Result struct {
	Target string
	Key    string
	Tags   []string
	Action string // "added", "removed", "unchanged"
}

// Options controls tag behaviour.
type Options struct {
	Target string
	Keys   []string // empty = all keys
	Add    []string
	Remove []string
}

// Apply adds or removes tags on variables in the named target.
func Apply(cfg *config.Document, opts Options) ([]Result, error) {
	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == opts.Target {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("target %q not found", opts.Target)
	}

	filter := toSet(opts.Keys)
	var results []Result

	for i := range target.Vars {
		v := &target.Vars[i]
		if len(filter) > 0 && !filter[v.Key] {
			continue
		}

		before := strings.Join(v.Tags, ",")
		v.Tags = applyChanges(v.Tags, opts.Add, opts.Remove)
		after := strings.Join(v.Tags, ",")

		action := "unchanged"
		if before != after {
			action = "updated"
		}

		results = append(results, Result{
			Target: target.Name,
			Key:    v.Key,
			Tags:   append([]string(nil), v.Tags...),
			Action: action,
		})
	}
	return results, nil
}

// List returns all unique tags used across a target's variables.
func List(cfg *config.Document, targetName string) ([]string, error) {
	for _, t := range cfg.Targets {
		if t.Name != targetName {
			continue
		}
		seen := map[string]struct{}{}
		for _, v := range t.Vars {
			for _, tag := range v.Tags {
				seen[tag] = struct{}{}
			}
		}
		out := make([]string, 0, len(seen))
		for tag := range seen {
			out = append(out, tag)
		}
		sort.Strings(out)
		return out, nil
	}
	return nil, fmt.Errorf("target %q not found", targetName)
}

func applyChanges(current, add, remove []string) []string {
	set := toSet(current)
	for _, t := range add {
		set[t] = true
	}
	for _, t := range remove {
		delete(set, t)
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
