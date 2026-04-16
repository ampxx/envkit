package envset_test

import (
	"testing"

	"envkit/internal/config"
	"envkit/internal/envset"
)

func makeConfig(targetName string, vars ...config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: targetName, Vars: vars},
		},
	}
}

func varDef(key, value string) config.VarDef {
	return config.VarDef{Key: key, Value: value}
}

func TestApply_SetsNewKey(t *testing.T) {
	cfg := makeConfig("production")
	results, err := envset.Apply(cfg, map[string]string{"FOO": "bar"}, envset.Options{Target: "production", Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Created {
		t.Errorf("expected Created=true, got %+v", results)
	}
	if cfg.Targets[0].Vars[0].Value != "bar" {
		t.Errorf("expected value bar, got %s", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig("production", varDef("FOO", "original"))
	results, err := envset.Apply(cfg, map[string]string{"FOO": "new"}, envset.Options{Target: "production", Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected Skipped=true, got %+v", results)
	}
	if cfg.Targets[0].Vars[0].Value != "original" {
		t.Errorf("value should not change, got %s", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_OverwritesExisting(t *testing.T) {
	cfg := makeConfig("staging", varDef("DB_URL", "old"))
	_, err := envset.Apply(cfg, map[string]string{"DB_URL": "new"}, envset.Options{Target: "staging", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets[0].Vars[0].Value != "new" {
		t.Errorf("expected updated value, got %s", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_DryRunDoesNotMutate(t *testing.T) {
	cfg := makeConfig("production")
	_, err := envset.Apply(cfg, map[string]string{"NEW_KEY": "val"}, envset.Options{Target: "production", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Targets[0].Vars) != 0 {
		t.Errorf("dry run should not mutate config")
	}
}

func TestApply_UnknownTargetErrors(t *testing.T) {
	cfg := makeConfig("production")
	_, err := envset.Apply(cfg, map[string]string{"K": "v"}, envset.Options{Target: "ghost"})
	if err == nil {
		t.Error("expected error for unknown target")
	}
}
