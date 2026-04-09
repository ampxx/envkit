package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestLog_CreatesFile(t *testing.T) {
	path := tempLog(t)
	l := New(path)

	if err := l.Log(EventValidate, "production", true, "all checks passed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected log file to be created")
	}
}

func TestLog_RoundTrip(t *testing.T) {
	path := tempLog(t)
	l := New(path)

	_ = l.Log(EventExport, "staging", true, "exported dotenv")
	_ = l.Log(EventLint, "staging", false, "2 issues found")

	entries, err := ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Event != EventExport {
		t.Errorf("expected event %q, got %q", EventExport, entries[0].Event)
	}
	if entries[1].Success {
		t.Error("expected second entry to have Success=false")
	}
}

func TestReadAll_EmptyFile(t *testing.T) {
	path := tempLog(t)
	entries, err := ReadAll(path)
	if err != nil {
		t.Fatalf("unexpected error on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestLog_MultipleEvents(t *testing.T) {
	path := tempLog(t)
	l := New(path)

	events := []EventType{EventValidate, EventFill, EventSnapshot}
	for _, ev := range events {
		if err := l.Log(ev, "dev", true, ""); err != nil {
			t.Fatalf("log error for %s: %v", ev, err)
		}
	}

	entries, err := ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if len(entries) != len(events) {
		t.Fatalf("expected %d entries, got %d", len(events), len(entries))
	}
	for i, e := range entries {
		if e.Event != events[i] {
			t.Errorf("entry %d: expected %q, got %q", i, events[i], e.Event)
		}
	}
}
