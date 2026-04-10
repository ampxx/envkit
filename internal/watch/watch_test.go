package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestHashFile_ReturnsConsistentHash(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, ".env", "KEY=value\n")
	h1, err := hashFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, _ := hashFile(p)
	if h1 != h2 {
		t.Errorf("expected same hash, got %q vs %q", h1, h2)
	}
}

func TestHashFile_ChangesOnWrite(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, ".env", "KEY=value\n")
	h1, _ := hashFile(p)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}
	h2, _ := hashFile(p)
	if h1 == h2 {
		t.Error("expected different hashes after file change")
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, ".env", "KEY=original\n")

	w := New([]string{p}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=updated\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != p {
			t.Errorf("expected path %q, got %q", p, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected old and new hashes to differ")
		}
	case <-time.After(300 * time.Millisecond):
		t.Error("timed out waiting for file change event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, ".env", "KEY=stable\n")

	w := New([]string{p}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// expected — no changes
	}
}
