package pin

import (
	"fmt"
	"time"

	"envkit/internal/config"
)

// Entry represents a pinned value for a variable in a target.
type Entry struct {
	Target    string    `yaml:"target"`
	Key       string    `yaml:"key"`
	Value     string    `yaml:"value"`
	PinnedAt  time.Time `yaml:"pinned_at"`
	PinnedBy  string    `yaml:"pinned_by"`
	ExpiresAt time.Time `yaml:"expires_at,omitempty"`
}

// Result holds the outcome of a pin operation.
type Result struct {
	Pinned  []Entry
	Skipped []string
	Errors  []string
}

// Pin locks one or more keys in a target to their current values.
// If keys is empty, all variables in the target are pinned.
func Pin(cfg *config.Document, targetName string, keys []string, pinnedBy string, ttl time.Duration) (Result, error) {
	target, ok := findTarget(cfg, targetName)
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", targetName)
	}

	keySet := toSet(keys)
	now := time.Now().UTC()
	var result Result

	for _, v := range target.Vars {
		if len(keySet) > 0 && !keySet[v.Key] {
			result.Skipped = append(result.Skipped, v.Key)
			continue
		}
		if v.Default == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: no value to pin", v.Key))
			continue
		}
		e := Entry{
			Target:   targetName,
			Key:      v.Key,
			Value:    v.Default,
			PinnedAt: now,
			PinnedBy: pinnedBy,
		}
		if ttl > 0 {
			e.ExpiresAt = now.Add(ttl)
		}
		result.Pinned = append(result.Pinned, e)
	}
	return result, nil
}

// Expired returns entries whose ExpiresAt is before now.
func Expired(entries []Entry) []Entry {
	now := time.Now().UTC()
	var out []Entry
	for _, e := range entries {
		if !e.ExpiresAt.IsZero() && e.ExpiresAt.Before(now) {
			out = append(out, e)
		}
	}
	return out
}

func findTarget(cfg *config.Document, name string) (config.Target, bool) {
	for _, t := range cfg.Targets {
		if t.Name == name {
			return t, true
		}
	}
	return config.Target{}, false
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
