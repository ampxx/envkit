package lint_test

import (
	"testing"

	"github.com/envkit/envkit/internal/config"
	"github.com/envkit/envkit/internal/lint"
)

func makeConfig(targets []config.Target) *config.Config {
	return &config.Config{Targets: targets}
}

func TestRun_NoIssues(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{
			Name: "production",
			Vars: []config.Var{
				{Key: "DATABASE_URL", Description: "Primary DB connection string", Required: true},
			},
		},
	})
	result := lint.Run(cfg)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
}

func TestRun_LowercaseKey(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{
			Name: "staging",
			Vars: []config.Var{
				{Key: "api_key", Description: "API key", Required: true},
			},
		},
	})
	result := lint.Run(cfg)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Level != "warn" {
		t.Errorf("expected warn, got %s", result.Issues[0].Level)
	}
}

func TestRun_MissingDescription(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{
			Name: "staging",
			Vars: []config.Var{
				{Key: "SECRET_KEY", Description: "", Required: true},
			},
		},
	})
	result := lint.Run(cfg)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Level != "warn" {
		t.Errorf("expected warn level, got %s", result.Issues[0].Level)
	}
}

func TestRun_DuplicateKey(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{
			Name: "production",
			Vars: []config.Var{
				{Key: "PORT", Description: "App port", Required: true},
				{Key: "PORT", Description: "App port duplicate", Required: false},
			},
		},
	})
	result := lint.Run(cfg)
	hasError := false
	for _, issue := range result.Issues {
		if issue.Level == "error" {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected an error-level issue for duplicate key")
	}
	if !result.HasErrors() {
		t.Error("HasErrors() should return true")
	}
}
