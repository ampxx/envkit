package extract

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
					{Key: "DATABASE_URL", Default: "postgres://prod/db"},
					{Key: "API_KEY", Default: "secret"},
					{Key: "LOG_LEVEL", Default: "warn"},
				},
			},
		},
	}
}

func TestApply_AllVars(t *testing.T) {
	cfg := makeConfig()
	results, err := Apply(cfg, Options{Target: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestApply_FilterByKey(t *testing.T) {
	cfg := makeConfig()
	results, err := Apply(cfg, Options{Target: "production", Keys: []string{"API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "API_KEY" {
		t.Fatalf("expected API_KEY only, got %+v", results)
	}
}

func TestApply_FilterByPattern(t *testing.T) {
	cfg := makeConfig()
	results, err := Apply(cfg, Options{Target: "production", Pattern: "^LOG_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "LOG_LEVEL" {
		t.Fatalf("expected LOG_LEVEL, got %+v", results)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig()
	_, err := Apply(cfg, Options{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	cfg := makeConfig()
	_, err := Apply(cfg, Options{Target: "production", Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}
