package coerce

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(vars []config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: "prod", Vars: vars},
		},
	}
}

func varDef(key, typ string) config.VarDef {
	return config.VarDef{Key: key, Type: typ}
}

func TestApply_BoolNormalises(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("DEBUG", "bool")})
	env := map[string]string{"DEBUG": "TRUE"}
	results, err := Apply(cfg, env, Options{Target: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Coerced != "true" {
		t.Errorf("expected coerced=true, got %+v", results)
	}
	if env["DEBUG"] != "true" {
		t.Errorf("env not updated, got %q", env["DEBUG"])
	}
}

func TestApply_IntTrimsSpace(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("PORT", "int")})
	env := map[string]string{"PORT": "  8080  "}
	results, err := Apply(cfg, env, Options{Target: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Coerced != "8080" {
		t.Errorf("expected 8080, got %q", results[0].Coerced)
	}
}

func TestApply_DryRunDoesNotMutate(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("DEBUG", "bool")})
	env := map[string]string{"DEBUG": "1"}
	_, err := Apply(cfg, env, Options{Target: "prod", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DEBUG"] != "1" {
		t.Errorf("dry-run mutated env: got %q", env["DEBUG"])
	}
}

func TestApply_SkipsMissingEnvKey(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("MISSING", "int")})
	env := map[string]string{}
	results, err := Apply(cfg, env, Options{Target: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Errorf("expected Skipped=true")
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig(nil)
	_, err := Apply(cfg, map[string]string{}, Options{Target: "nope"})
	if err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestApply_FilterByKeys(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		varDef("PORT", "int"),
		varDef("DEBUG", "bool"),
	})
	env := map[string]string{"PORT": " 9000 ", "DEBUG": "TRUE"}
	results, err := Apply(cfg, env, Options{Target: "prod", Keys: []string{"PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PORT" {
		t.Errorf("expected only PORT in results, got %+v", results)
	}
}
