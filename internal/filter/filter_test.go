package filter

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

func varDef(key, val string, tags ...string) config.VarDef {
	return config.VarDef{Key: key, Default: val, Tags: tags}
}

func TestApply_ReturnsAllVars(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DB_HOST", "localhost"),
		varDef("API_KEY", "secret"),
	})
	results, err := Apply(cfg, Options{Target: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestApply_FiltersByPrefix(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DB_HOST", "localhost"),
		varDef("DB_PORT", "5432"),
		varDef("API_KEY", "secret"),
	})
	results, err := Apply(cfg, Options{Target: "production", Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "API_KEY" {
			t.Errorf("API_KEY should have been filtered out")
		}
	}
}

func TestApply_FiltersByPattern(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DB_HOST", "localhost"),
		varDef("REDIS_HOST", "redis"),
		varDef("API_KEY", "secret"),
	})
	results, err := Apply(cfg, Options{Target: "production", Pattern: "HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestApply_FiltersByTags(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DB_HOST", "localhost", "db", "infra"),
		varDef("API_KEY", "secret", "auth"),
		varDef("DB_PORT", "5432", "db"),
	})
	results, err := Apply(cfg, Options{Target: "production", Tags: []string{"db", "infra"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Errorf("expected only DB_HOST, got %+v", results)
	}
}

func TestApply_UnknownTargetErrors(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("X", "y")})
	_, err := Apply(cfg, Options{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_KeyAllowlist(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DB_HOST", "localhost"),
		varDef("API_KEY", "secret"),
		varDef("LOG_LEVEL", "info"),
	})
	results, err := Apply(cfg, Options{Target: "production", Keys: []string{"API_KEY", "LOG_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
