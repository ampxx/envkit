package summary

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig() *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{
				Name: "production",
				Vars: []config.VarDef{
					{Key: "DB_URL", Required: true, Sensitive: true, Tags: []string{"db"}},
					{Key: "PORT", Required: true, Default: "8080"},
					{Key: "LOG_LEVEL", Required: false, Default: "info"},
				},
			},
			{
				Name: "staging",
				Vars: []config.VarDef{
					{Key: "DB_URL", Required: true, Sensitive: true},
					{Key: "DEBUG", Required: false},
				},
			},
		},
	}
}

func TestRun_AllTargets(t *testing.T) {
	cfg := makeConfig()
	report, err := Run(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(report.Results))
	}
}

func TestRun_FilterByTarget(t *testing.T) {
	cfg := makeConfig()
	report, err := Run(cfg, []string{"staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(report.Results))
	}
	if report.Results[0].Target != "staging" {
		t.Errorf("expected staging, got %s", report.Results[0].Target)
	}
}

func TestRun_ProductionCounts(t *testing.T) {
	cfg := makeConfig()
	report, err := Run(cfg, []string{"production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := report.Results[0]
	if r.TotalVars != 3 {
		t.Errorf("TotalVars: want 3, got %d", r.TotalVars)
	}
	if r.Required != 2 {
		t.Errorf("Required: want 2, got %d", r.Required)
	}
	if r.Optional != 1 {
		t.Errorf("Optional: want 1, got %d", r.Optional)
	}
	if r.Sensitive != 1 {
		t.Errorf("Sensitive: want 1, got %d", r.Sensitive)
	}
	if r.Tagged != 1 {
		t.Errorf("Tagged: want 1, got %d", r.Tagged)
	}
	if r.WithDefault != 2 {
		t.Errorf("WithDefault: want 2, got %d", r.WithDefault)
	}
}

func TestRun_UnknownTarget(t *testing.T) {
	cfg := makeConfig()
	_, err := Run(cfg, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown target, got nil")
	}
}

func TestRun_SortedAlphabetically(t *testing.T) {
	cfg := makeConfig()
	report, err := Run(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Results[0].Target != "production" {
		t.Errorf("expected production first, got %s", report.Results[0].Target)
	}
	if report.Results[1].Target != "staging" {
		t.Errorf("expected staging second, got %s", report.Results[1].Target)
	}
}
