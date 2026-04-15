package inherit

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func varDef(key, value string) config.VarDef {
	return config.VarDef{Key: key, Value: value}
}

func makeTarget(name string, vars ...config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func TestApply_InheritsNewKeys(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_HOST", "staging-db"), varDef("PORT", "5432")),
		makeTarget("production"),
	)
	results, err := Apply(cfg, "staging", "production", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Skipped {
			t.Errorf("expected %q not to be skipped", r.Key)
		}
	}
	prod := findTarget(cfg, "production")
	if len(prod.Vars) != 2 {
		t.Errorf("expected 2 vars in production, got %d", len(prod.Vars))
	}
}

func TestApply_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_HOST", "staging-db")),
		makeTarget("production", varDef("DB_HOST", "prod-db")),
	)
	results, err := Apply(cfg, "staging", "production", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected DB_HOST to be skipped")
	}
	prod := findTarget(cfg, "production")
	if prod.Vars[0].Value != "prod-db" {
		t.Errorf("expected prod-db, got %s", prod.Vars[0].Value)
	}
}

func TestApply_OverwritesExisting(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_HOST", "staging-db")),
		makeTarget("production", varDef("DB_HOST", "prod-db")),
	)
	_, err := Apply(cfg, "staging", "production", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := findTarget(cfg, "production")
	if prod.Vars[0].Value != "staging-db" {
		t.Errorf("expected staging-db after overwrite, got %s", prod.Vars[0].Value)
	}
}

func TestApply_FiltersByKeys(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", varDef("DB_HOST", "staging-db"), varDef("PORT", "5432")),
		makeTarget("production"),
	)
	results, err := Apply(cfg, "staging", "production", Options{Keys: []string{"PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PORT" {
		t.Errorf("expected only PORT to be inherited")
	}
}

func TestApply_UnknownSourceErrors(t *testing.T) {
	cfg := makeConfig(makeTarget("production"))
	_, err := Apply(cfg, "staging", "production", Options{})
	if err == nil {
		t.Fatal("expected error for unknown source target")
	}
}

func TestApply_UnknownDestErrors(t *testing.T) {
	cfg := makeConfig(makeTarget("staging"))
	_, err := Apply(cfg, "staging", "production", Options{})
	if err == nil {
		t.Fatal("expected error for unknown destination target")
	}
}
