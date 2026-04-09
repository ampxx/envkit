package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	path, err := Save("prod", vars, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected snapshot file to exist at %s", path)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	path, err := Save("staging", vars, dir)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	snap, err := Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if snap.Target != "staging" {
		t.Errorf("expected target 'staging', got %q", snap.Target)
	}

	for k, v := range vars {
		if snap.Vars[k] != v {
			t.Errorf("expected vars[%q]=%q, got %q", k, v, snap.Vars[k])
		}
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(badPath, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(badPath)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestList_ReturnsMatchingFiles(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{"KEY": "val"}

	if _, err := Save("dev", vars, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := Save("dev", vars, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := Save("prod", vars, dir); err != nil {
		t.Fatal(err)
	}

	devSnaps, err := List("dev", dir)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(devSnaps) != 2 {
		t.Errorf("expected 2 dev snapshots, got %d", len(devSnaps))
	}

	prodSnaps, err := List("prod", dir)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(prodSnaps) != 1 {
		t.Errorf("expected 1 prod snapshot, got %d", len(prodSnaps))
	}
}
