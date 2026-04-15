package strict

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(vars []config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: "production", Vars: vars},
		},
	}
}

func varDef(key, value, desc, def string) config.VarDef {
	return config.VarDef{Key: key, Value: value, Description: desc, Default: def}
}

func TestRun_NoIssues(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DB_HOST", "localhost", "database host", ""),
		varDef("PORT", "5432", "server port", ""),
	})
	report, err := Run(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasFailures() {
		t.Errorf("expected no failures, got %d", len(report.Results))
	}
}

func TestRun_LowercaseKey(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("db_host", "localhost", "database host", ""),
	})
	report, err := Run(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasFailures() {
		t.Fatal("expected failure for lowercase key")
	}
	if report.Results[0].Issue != "key contains lowercase letters" {
		t.Errorf("unexpected issue: %s", report.Results[0].Issue)
	}
}

func TestRun_MissingDescription(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("API_KEY", "secret", "", ""),
	})
	report, err := Run(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range report.Results {
		if r.Issue == "missing description" {
			found = true
		}
	}
	if !found {
		t.Error("expected missing description issue")
	}
}

func TestRun_EmptyValueNoDefault(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("SECRET_KEY", "", "some secret", ""),
	})
	report, err := Run(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range report.Results {
		if r.Issue == "empty value with no default" {
			found = true
		}
	}
	if !found {
		t.Error("expected empty value issue")
	}
}

func TestRun_EmptyValueWithDefaultIsOK(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("LOG_LEVEL", "", "logging level", "info"),
	})
	report, err := Run(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range report.Results {
		if r.Issue == "empty value with no default" {
			t.Errorf("should not flag empty value when default is set")
		}
	}
}

func TestRun_UnknownTarget(t *testing.T) {
	cfg := makeConfig(nil)
	_, err := Run(cfg, "staging")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}
