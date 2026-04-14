package dedupe

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

func TestApply_NoDuplicates(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "production",
		Vars: []config.VarDef{varDef("A", "1"), varDef("B", "2")},
	})
	report, err := Apply(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Total != 0 {
		t.Errorf("expected 0 removals, got %d", report.Total)
	}
	if len(cfg.Targets[0].Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(cfg.Targets[0].Vars))
	}
}

func TestApply_RemovesDuplicates(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "staging",
		Vars: []config.VarDef{
			varDef("DB_HOST", "old"),
			varDef("DB_HOST", "new"),
			varDef("PORT", "8080"),
		},
	})
	report, err := Apply(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Total != 1 {
		t.Errorf("expected 1 removal, got %d", report.Total)
	}
	if len(cfg.Targets[0].Vars) != 2 {
		t.Errorf("expected 2 vars after dedup, got %d", len(cfg.Targets[0].Vars))
	}
	// last-write-wins: kept value should be "new"
	if cfg.Targets[0].Vars[0].Value != "new" {
		t.Errorf("expected kept value 'new', got %q", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_FiltersByTarget(t *testing.T) {
	cfg := makeConfig(
		config.Target{
			Name: "dev",
			Vars: []config.VarDef{varDef("X", "1"), varDef("X", "2")},
		},
		config.Target{
			Name: "prod",
			Vars: []config.VarDef{varDef("X", "a"), varDef("X", "b")},
		},
	)
	report, err := Apply(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Total != 1 {
		t.Errorf("expected 1 removal in dev, got %d", report.Total)
	}
	// prod target should be untouched
	if len(cfg.Targets[1].Vars) != 2 {
		t.Errorf("prod target should still have 2 vars")
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig(config.Target{Name: "dev", Vars: []config.VarDef{varDef("A", "1")}})
	_, err := Apply(cfg, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown target, got nil")
	}
}
