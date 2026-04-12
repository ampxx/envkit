package copy

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targets ...config.Target) config.Document {
	return config.Document{Targets: targets}
}

func varDef(key, value string) config.VarDef {
	return config.VarDef{Key: key, Default: value}
}

func makeTarget(name string, vars ...config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func TestCopy_CopiesAllVars(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_URL", "postgres://staging"), varDef("PORT", "5432")),
		makeTarget("prod"),
	)
	results, updated, err := Copy(cfg, "staging", "prod", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	prod := findTarget(updated, "prod")
	if prod == nil || len(prod.Vars) != 2 {
		t.Errorf("expected 2 vars in prod, got %v", prod)
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_URL", "postgres://staging")),
		makeTarget("prod", varDef("DB_URL", "postgres://prod")),
	)
	results, updated, err := Copy(cfg, "staging", "prod", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected 1 skipped result, got %+v", results)
	}
	prod := findTarget(updated, "prod")
	if prod.Vars[0].Default != "postgres://prod" {
		t.Errorf("expected original value preserved")
	}
}

func TestCopy_OverwritesExisting(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_URL", "postgres://staging")),
		makeTarget("prod", varDef("DB_URL", "postgres://prod")),
	)
	_, updated, err := Copy(cfg, "staging", "prod", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := findTarget(updated, "prod")
	if prod.Vars[0].Default != "postgres://staging" {
		t.Errorf("expected overwritten value, got %s", prod.Vars[0].Default)
	}
}

func TestCopy_FiltersByKey(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_URL", "postgres://staging"), varDef("PORT", "5432")),
		makeTarget("prod"),
	)
	results, updated, err := Copy(cfg, "staging", "prod", Options{Keys: []string{"PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PORT" {
		t.Errorf("expected only PORT copied, got %+v", results)
	}
	prod := findTarget(updated, "prod")
	if len(prod.Vars) != 1 || prod.Vars[0].Key != "PORT" {
		t.Errorf("expected prod to have only PORT")
	}
}

func TestCopy_MissingSourceErrors(t *testing.T) {
	cfg := makeConfig(makeTarget("prod"))
	_, _, err := Copy(cfg, "staging", "prod", Options{})
	if err == nil {
		t.Error("expected error for missing source target")
	}
}

func TestCopy_MissingDestErrors(t *testing.T) {
	cfg := makeConfig(makeTarget("staging", varDef("PORT", "5432")))
	_, _, err := Copy(cfg, "staging", "prod", Options{})
	if err == nil {
		t.Error("expected error for missing destination target")
	}
}
