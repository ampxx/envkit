package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvContent(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestValidateFile_NoIssues(t *testing.T) {
	p := writeTempEnvContent(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	errs, err := ValidateFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateFile_LowercaseKey(t *testing.T) {
	p := writeTempEnvContent(t, "app_host=localhost\n")
	errs, err := ValidateFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected error for lowercase key")
	}
}

func TestValidateFile_DuplicateKey(t *testing.T) {
	p := writeTempEnvContent(t, "APP_HOST=localhost\nAPP_HOST=remotehost\n")
	errs, err := ValidateFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, e := range errs {
		if e.Key == "APP_HOST" && e.Line == 2 {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate key error for APP_HOST on line 2, got %v", errs)
	}
}

func TestValidateFile_MissingFile(t *testing.T) {
	_, err := ValidateFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidationSummary(t *testing.T) {
	if got := ValidationSummary(nil); got != "no issues found" {
		t.Errorf("unexpected summary: %q", got)
	}
	errs := []ValidationError{{Line: 1, Key: "x", Message: "bad"}}
	if got := ValidationSummary(errs); got != "1 issue(s) found" {
		t.Errorf("unexpected summary: %q", got)
	}
}
