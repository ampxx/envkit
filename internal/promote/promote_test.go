package promote

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig() *config.Config {
	return &config.Config{
		Targets: []config.Target{
			{
				Name: "staging",
				Vars: []config.VarDef{
					{Key: "API_URL", Default: "https://staging.example.com"},
					{Key: "DEBUG", Default: "true"},
				},
			},
			{
				Name: "production",
				Vars: []config.VarDef{
					{Key: "API_URL", Default: "https://prod.example.com"},
				},
			},
		},
	}
}

func TestPromote_CopiesNewKey(t *testing.T) {
	cfg := makeConfig()
	results, err := Promote(cfg, "staging", "production", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// DEBUG should be added to production
	dst := findTarget(cfg, "production")
	if findVar(dst, "DEBUG") == nil {
		t.Error("expected DEBUG to be promoted to production")
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig()
	results, err := Promote(cfg, "staging", "production", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var skipped *Result
	for i := range results {
		if results[i].Key == "API_URL" {
			skipped = &results[i]
		}
	}
	if skipped == nil || !skipped.Skipped {
		t.Error("expected API_URL to be skipped")
	}
}

func TestPromote_OverwritesExisting(t *testing.T) {
	cfg := makeConfig()
	_, err := Promote(cfg, "staging", "production", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dst := findTarget(cfg, "production")
	v := findVar(dst, "API_URL")
	if v == nil || v.Default != "https://staging.example.com" {
		t.Errorf("expected API_URL overwritten, got %v", v)
	}
}

func TestPromote_FiltersByKey(t *testing.T) {
	cfg := makeConfig()
	results, err := Promote(cfg, "staging", "production", Options{Keys: []string{"DEBUG"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DEBUG" {
		t.Errorf("expected only DEBUG in results, got %v", results)
	}
}

func TestPromote_InvalidSource(t *testing.T) {
	cfg := makeConfig()
	_, err := Promote(cfg, "nonexistent", "production", Options{})
	if err == nil {
		t.Error("expected error for missing source target")
	}
}

func TestPromote_InvalidDest(t *testing.T) {
	cfg := makeConfig()
	_, err := Promote(cfg, "staging", "nonexistent", Options{})
	if err == nil {
		t.Error("expected error for missing destination target")
	}
}
