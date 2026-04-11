package defaults

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targetName string, vars []config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: targetName, Vars: vars},
		},
	}
}

func varDef(key, value, def string) config.VarDef {
	return config.VarDef{Key: key, Value: value, Default: def}
}

func TestApply_FillsEmptyValue(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("PORT", "", "8080"),
	})
	results, err := Apply(cfg, Options{Target: "dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Applied {
		t.Fatalf("expected PORT to be applied")
	}
	if cfg.Targets[0].Vars[0].Value != "8080" {
		t.Errorf("expected value 8080, got %q", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_SkipsNonEmptyWithoutOverwrite(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("PORT", "3000", "8080"),
	})
	results, err := Apply(cfg, Options{Target: "dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Applied {
		t.Fatalf("expected PORT to be skipped")
	}
	if cfg.Targets[0].Vars[0].Value != "3000" {
		t.Errorf("value should remain 3000, got %q", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_OverwriteExisting(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("PORT", "3000", "8080"),
	})
	results, err := Apply(cfg, Options{Target: "dev", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Fatalf("expected PORT to be overwritten")
	}
	if cfg.Targets[0].Vars[0].Value != "8080" {
		t.Errorf("expected value 8080, got %q", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_FiltersByKey(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("PORT", "", "8080"),
		varDef("HOST", "", "localhost"),
	})
	results, err := Apply(cfg, Options{Target: "dev", Keys: []string{"PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PORT" {
		t.Fatalf("expected only PORT in results")
	}
}

func TestApply_UnknownTargetErrors(t *testing.T) {
	cfg := makeConfig("dev", nil)
	_, err := Apply(cfg, Options{Target: "prod"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_SkipsVarsWithNoDefault(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("SECRET", "", ""),
	})
	results, err := Apply(cfg, Options{Target: "dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results for var with no default")
	}
}
