package rename

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func varDef(key, value string) config.VarDef {
	return config.VarDef{Key: key, Default: value}
}

func TestApply_RenamesKey(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "production",
		Vars: []config.VarDef{varDef("OLD_KEY", "val")},
	})

	results, err := Apply(cfg, Options{Keys: map[string]string{"OLD_KEY": "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Renamed {
		t.Fatalf("expected rename, got %+v", results)
	}
	if cfg.Targets[0].Vars[0].Key != "NEW_KEY" {
		t.Errorf("expected key NEW_KEY, got %s", cfg.Targets[0].Vars[0].Key)
	}
}

func TestApply_KeyNotFound(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "staging",
		Vars: []config.VarDef{varDef("EXISTING", "v")},
	})

	results, _ := Apply(cfg, Options{Keys: map[string]string{"MISSING": "NEW"}})
	if results[0].Renamed {
		t.Error("expected not renamed")
	}
	if results[0].Reason != "key not found" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestApply_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "dev",
		Vars: []config.VarDef{varDef("A", "1"), varDef("B", "2")},
	})

	results, _ := Apply(cfg, Options{Keys: map[string]string{"A": "B"}, Overwrite: false})
	if results[0].Renamed {
		t.Error("expected skip when destination exists and overwrite=false")
	}
}

func TestApply_OverwritesExisting(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "dev",
		Vars: []config.VarDef{varDef("A", "1"), varDef("B", "2")},
	})

	results, _ := Apply(cfg, Options{Keys: map[string]string{"A": "B"}, Overwrite: true})
	if !results[0].Renamed {
		t.Errorf("expected rename with overwrite, got: %+v", results[0])
	}
	if len(cfg.Targets[0].Vars) != 1 {
		t.Errorf("expected 1 var after overwrite, got %d", len(cfg.Targets[0].Vars))
	}
}

func TestApply_FiltersByTarget(t *testing.T) {
	cfg := makeConfig(
		config.Target{Name: "dev", Vars: []config.VarDef{varDef("KEY", "v")}},
		config.Target{Name: "prod", Vars: []config.VarDef{varDef("KEY", "v")}},
	)

	_, _ = Apply(cfg, Options{Target: "dev", Keys: map[string]string{"KEY": "NEW_KEY"}})

	if cfg.Targets[0].Vars[0].Key != "NEW_KEY" {
		t.Errorf("dev target should have NEW_KEY")
	}
	if cfg.Targets[1].Vars[0].Key != "KEY" {
		t.Errorf("prod target should be unchanged")
	}
}

func TestApply_NoMappingsError(t *testing.T) {
	cfg := makeConfig(config.Target{Name: "dev"})
	_, err := Apply(cfg, Options{})
	if err == nil {
		t.Error("expected error for empty key mappings")
	}
}
