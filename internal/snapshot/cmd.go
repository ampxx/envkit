package snapshot

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envkit/internal/config"
)

const defaultSnapshotDir = ".envkit/snapshots"

// RunSnapshot saves a snapshot of the resolved env vars for a given target.
func RunSnapshot(configPath, target, snapshotDir string) error {
	if snapshotDir == "" {
		snapshotDir = defaultSnapshotDir
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	var tgt *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == target {
			tgt = &cfg.Targets[i]
			break
		}
	}
	if tgt == nil {
		return fmt.Errorf("target %q not found in config", target)
	}

	vars := make(map[string]string, len(tgt.Vars))
	for k, v := range tgt.Vars {
		vars[k] = v
	}

	path, err := Save(target, vars, snapshotDir)
	if err != nil {
		return fmt.Errorf("saving snapshot: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot saved: %s\n", filepath.Clean(path))
	return nil
}

// RunListSnapshots prints all saved snapshots for a given target.
func RunListSnapshots(target, snapshotDir string) error {
	if snapshotDir == "" {
		snapshotDir = defaultSnapshotDir
	}

	paths, err := List(target, snapshotDir)
	if err != nil {
		return fmt.Errorf("listing snapshots: %w", err)
	}

	if len(paths) == 0 {
		fmt.Fprintf(os.Stdout, "No snapshots found for target %q\n", target)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Snapshots for target %q:\n", target)
	for _, p := range paths {
		snap, err := Load(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  [error reading %s: %v]\n", p, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "  %s  (%d vars)\n", snap.Timestamp.Format("2006-01-02 15:04:05 UTC"), len(snap.Vars))
	}
	return nil
}
