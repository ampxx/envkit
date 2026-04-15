package envcount_test

import (
	"testing"

	"envkit/internal/config"
	"envkit/internal/envcount"
)

func makeConfig() *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{
				Name: "production",
				Vars: []config.VarDef{
					{Key: "DB_URL", Required: true, Sensitive: true},
					{Key: "PORT", Required: true},
					{Key: "LOG_LEVEL"},
				},
			},
			{
				Name: "staging",
				Vars: []config.VarDef{
					{Key: "DB_URL", Required: true, Sensitive: true},
					{Key: "DEBUG"},
				},
			},
		},
	}
}

func TestCount_AllTargets(t *testing.T) {
	cfg := makeConfig()
	report, err := envcount.Count(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(report.Results))
	}
	if report.GrandTotal() != 5 {
		t.Errorf("expected grand total 5, got %d", report.GrandTotal())
	}
}

func TestCount_FilterByTarget(t *testing.T) {
	cfg := makeConfig()
	report, err := envcount.Count(cfg, []string{"production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(report.Results))
	}
	r := report.Results[0]
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
	if r.Required != 2 {
		t.Errorf("expected Required=2, got %d", r.Required)
	}
	if r.Optional != 1 {
		t.Errorf("expected Optional=1, got %d", r.Optional)
	}
	if r.Sensitive != 1 {
		t.Errorf("expected Sensitive=1, got %d", r.Sensitive)
	}
}

func TestCount_UnknownTarget(t *testing.T) {
	cfg := makeConfig()
	_, err := envcount.Count(cfg, []string{"nope"})
	if err == nil {
		t.Fatal("expected error for unknown target, got nil")
	}
}

func TestCount_ResultsSortedByName(t *testing.T) {
	cfg := makeConfig()
	report, err := envcount.Count(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Results[0].Target != "production" {
		t.Errorf("expected first target to be production, got %s", report.Results[0].Target)
	}
	if report.Results[1].Target != "staging" {
		t.Errorf("expected second target to be staging, got %s", report.Results[1].Target)
	}
}

func TestGrandTotal_Empty(t *testing.T) {
	r := envcount.Report{}
	if r.GrandTotal() != 0 {
		t.Errorf("expected 0 for empty report, got %d", r.GrandTotal())
	}
}
