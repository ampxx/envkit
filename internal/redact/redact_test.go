package redact_test

import (
	"testing"

	"github.com/yourusername/envkit/internal/redact"
)

func varsMap(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func findResult(results []redact.Result, key string) (redact.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return redact.Result{}, false
}

func TestApply_RedactsSensitiveKey(t *testing.T) {
	vars := varsMap("DB_PASSWORD", "s3cr3t", "APP_NAME", "envkit")
	results := redact.Apply(vars, nil)

	r, ok := findResult(results, "DB_PASSWORD")
	if !ok {
		t.Fatal("expected DB_PASSWORD in results")
	}
	if !r.Redacted {
		t.Error("expected DB_PASSWORD to be marked redacted")
	}
	if r.Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Value)
	}
}

func TestApply_PreservesNonSensitiveKey(t *testing.T) {
	vars := varsMap("APP_NAME", "envkit")
	results := redact.Apply(vars, nil)

	r, ok := findResult(results, "APP_NAME")
	if !ok {
		t.Fatal("expected APP_NAME in results")
	}
	if r.Redacted {
		t.Error("expected APP_NAME not to be redacted")
	}
	if r.Value != "envkit" {
		t.Errorf("expected envkit, got %q", r.Value)
	}
}

func TestApply_CustomRules(t *testing.T) {
	import_re := redact.Rule{
		Pattern:     mustCompile(`(?i)internal`),
		Replacement: "[HIDDEN]",
	}
	vars := varsMap("INTERNAL_HOST", "10.0.0.1", "PUBLIC_URL", "https://example.com")
	results := redact.Apply(vars, []redact.Rule{import_re})

	r, ok := findResult(results, "INTERNAL_HOST")
	if !ok {
		t.Fatal("expected INTERNAL_HOST in results")
	}
	if r.Value != "[HIDDEN]" {
		t.Errorf("expected [HIDDEN], got %q", r.Value)
	}
}

func TestApply_TokenKey(t *testing.T) {
	vars := varsMap("GITHUB_TOKEN", "ghp_abc123")
	results := redact.Apply(vars, nil)

	r, ok := findResult(results, "GITHUB_TOKEN")
	if !ok {
		t.Fatal("expected GITHUB_TOKEN in results")
	}
	if !r.Redacted {
		t.Error("expected GITHUB_TOKEN to be redacted")
	}
}

func mustCompile(s string) *regexp.Regexp {
	return regexp.MustCompile(s)
}
