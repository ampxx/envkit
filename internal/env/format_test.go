package env

import (
	"strings"
	"testing"
)

func TestFormat_Default(t *testing.T) {
	vars := map[string]string{"FOO": "bar"}
	out := Format(vars, FormatOption{})
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got %q", out)
	}
}

func TestFormat_Sorted(t *testing.T) {
	vars := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	out := Format(vars, FormatOption{Sorted: true})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "AAA=") {
		t.Errorf("expected first line to start with AAA=, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZZZ=") {
		t.Errorf("expected last line to start with ZZZ=, got %q", lines[2])
	}
}

func TestFormat_Exported(t *testing.T) {
	vars := map[string]string{"PORT": "8080"}
	out := Format(vars, FormatOption{Exported: true})
	if !strings.Contains(out, "export PORT=8080") {
		t.Errorf("expected export prefix, got %q", out)
	}
}

func TestFormat_Quoted(t *testing.T) {
	vars := map[string]string{"MSG": "hello world"}
	out := Format(vars, FormatOption{Quote: true})
	// %q wraps in double quotes
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got %q", out)
	}
}

func TestFormat_ExportedAndQuoted(t *testing.T) {
	vars := map[string]string{"SECRET": "abc"}
	out := Format(vars, FormatOption{Exported: true, Quote: true})
	if !strings.Contains(out, `export SECRET="abc"`) {
		t.Errorf("expected export + quoted, got %q", out)
	}
}

func TestRedact_ReplacesAllValues(t *testing.T) {
	vars := map[string]string{"A": "secret", "B": "also-secret"}
	redacted := Redact(vars)
	for k, v := range redacted {
		if v != "***" {
			t.Errorf("key %s: expected ***, got %q", k, v)
		}
	}
	// original must not be mutated
	if vars["A"] != "secret" {
		t.Error("original map was mutated")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	out := Redact(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
