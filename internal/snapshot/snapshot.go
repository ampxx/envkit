package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a saved state of environment variables at a point in time.
type Snapshot struct {
	Target    string            `json:"target"`
	Timestamp time.Time         `json:"timestamp"`
	Vars      map[string]string `json:"vars"`
}

// Save writes a snapshot of the given env vars to the snapshots directory.
func Save(target string, vars map[string]string, dir string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating snapshot dir: %w", err)
	}

	snap := Snapshot{
		Target:    target,
		Timestamp: time.Now().UTC(),
		Vars:      vars,
	}

	filename := fmt.Sprintf("%s_%s.json", target, snap.Timestamp.Format("20060102T150405Z"))
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("writing snapshot file: %w", err)
	}

	return path, nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("parsing snapshot: %w", err)
	}

	return &snap, nil
}

// List returns all snapshot file paths for a given target in the directory.
func List(target, dir string) ([]string, error) {
	pattern := filepath.Join(dir, target+"_*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("listing snapshots: %w", err)
	}
	return matches, nil
}

// Latest returns the most recent snapshot for a given target in the directory,
// or nil if no snapshots exist. Because snapshot filenames embed a sortable
// timestamp, the lexicographically last entry is also the most recent.
func Latest(target, dir string) (*Snapshot, error) {
	matches, err := List(target, dir)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, nil
	}
	return Load(matches[len(matches)-1])
}
