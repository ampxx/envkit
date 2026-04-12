package required_test

import (
	"testing"

	"envkit/internal/config"
	"envkit/internal/required"
)

func makeConfig(vars []config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: "production", Vars: vars},
		},
	}
}

func varDef(key string, req bool, def string) config.VarDef {
	return config.VarDef{Key: key, Required: req, Default: def}
}

func TestCheck_AllPresent(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DATABASE_URL", true, ""),
		varDef("API_KEY", true, ""),
	})
	env := map[string]string{"DATABASE_URL": "postgres://", "API_KEY": "secret"}
	rep, err := required.Check(cfg, "production", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.HasFailures() {
		t.Error("expected no failures")
	}
	if len(rep.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(rep.Results))
	}
}

func TestCheck_MissingRequired(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("DATABASE_URL", true, ""),
		varDef("API_KEY", true, ""),
	})
	env := map[string]string{"DATABASE_URL": "postgres://"}
	rep, err := required.Check(cfg, "production", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rep.HasFailures() {
		t.Error("expected failures")
	}
}

func TestCheck_MissingButHasDefault(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("TIMEOUT", true, "30s"),
	})
	env := map[string]string{}
	rep, err := required.Check(cfg, "production", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.HasFailures() {
		t.Error("variable with default should not be a failure")
	}
}

func TestCheck_UnknownTarget(t *testing.T) {
	cfg := makeConfig(nil)
	_, err := required.Check(cfg, "staging", map[string]string{})
	if err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestCheck_IgnoresOptionalVars(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("OPTIONAL", false, ""),
		varDef("REQUIRED", true, ""),
	})
	env := map[string]string{"REQUIRED": "value"}
	rep, err := required.Check(cfg, "production", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 {
		t.Errorf("expected 1 result (required only), got %d", len(rep.Results))
	}
}
