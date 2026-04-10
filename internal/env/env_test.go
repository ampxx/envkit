package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestParseFile_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "APP_ENV" || entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestParseFile_IgnoresComments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nKEY=value\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestParseFile_StripsQuotes(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "my secret value" {
		t.Errorf("expected stripped value, got %q", entries[0].Value)
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestToMap(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	m := ToMap(entries)
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	original := []Entry{{Key: "FOO", Value: "bar"}, {Key: "NUM", Value: "42"}}

	if err := WriteFile(path, original); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	loaded, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(loaded) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(loaded))
	}
	for i, e := range original {
		if loaded[i].Key != e.Key || loaded[i].Value != e.Value {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, loaded[i], e)
		}
	}
}
