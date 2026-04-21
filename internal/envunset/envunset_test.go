package envunset

import (
	"testing"

	"github.com/your-org/envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func varDef(key, value string) config.VarDef {
	return config.VarDef{Key: key, Value: value}
}

func TestApply_RemovesKey(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "production",
		Vars: []config.VarDef{varDef("DB_URL", "postgres://"), varDef("PORT", "5432")},
	})
	_, err := Apply(cfg, Options{Keys: []string{"DB_URL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range cfg.Targets[0].Vars {
		if v.Key == "DB_URL" {
			t.Error("expected DB_URL to be removed")
		}
	}
}

func TestApply_KeyNotFound(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "staging",
		Vars: []config.VarDef{varDef("PORT", "3000")},
	})
	results, err := Apply(cfg, Options{Keys: []string{"MISSING_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || results[0].Found {
		t.Error("expected Found=false for missing key")
	}
}

func TestApply_DryRunDoesNotMutate(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "dev",
		Vars: []config.VarDef{varDef("SECRET", "abc"), varDef("PORT", "8080")},
	})
	_, err := Apply(cfg, Options{Keys: []string{"SECRET"}, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsKey(cfg.Targets[0].Vars, "SECRET") {
		t.Error("dry-run should not remove SECRET")
	}
}

func TestApply_FiltersByTarget(t *testing.T) {
	cfg := makeConfig(
		config.Target{Name: "dev", Vars: []config.VarDef{varDef("TOKEN", "dev-tok")}},
		config.Target{Name: "prod", Vars: []config.VarDef{varDef("TOKEN", "prod-tok")}},
	)
	_, err := Apply(cfg, Options{Keys: []string{"TOKEN"}, Target: "dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if containsKey(cfg.Targets[0].Vars, "TOKEN") {
		t.Error("TOKEN should be removed from dev")
	}
	if !containsKey(cfg.Targets[1].Vars, "TOKEN") {
		t.Error("TOKEN should remain in prod")
	}
}

func TestApply_EmptyKeysErrors(t *testing.T) {
	cfg := makeConfig(config.Target{Name: "dev"})
	_, err := Apply(cfg, Options{Keys: []string{}})
	if err == nil {
		t.Error("expected error when no keys provided")
	}
}

func TestApply_ResultFound(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "prod",
		Vars: []config.VarDef{varDef("API_KEY", "secret")},
	})
	results, _ := Apply(cfg, Options{Keys: []string{"API_KEY"}})
	if len(results) == 0 || !results[0].Found {
		t.Error("expected Found=true for existing key")
	}
}
