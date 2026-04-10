package clone

import (
	"testing"

	"github.com/envkit/envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Config {
	return &config.Config{Targets: targets}
}

func varDef(key, value string) config.VarDef {
	return config.VarDef{Key: key, Default: value}
}

func TestClone_CopiesAllVars(t *testing.T) {
	cfg := makeConfig(
		config.Target{Name: "staging", Vars: []config.VarDef{varDef("FOO", "bar"), varDef("BAZ", "qux")}},
	)
	res, err := Clone(cfg, "staging", "production", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 2 {
		t.Errorf("expected 2 copied, got %d", res.Copied)
	}
	if len(cfg.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(cfg.Targets))
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig(
		config.Target{Name: "staging", Vars: []config.VarDef{varDef("FOO", "new")}},
		config.Target{Name: "production", Vars: []config.VarDef{varDef("FOO", "old")}},
	)
	res, err := Clone(cfg, "staging", "production", Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
	prod := findTarget(cfg, "production")
	if prod.Vars[0].Default != "old" {
		t.Errorf("expected value to remain 'old', got %q", prod.Vars[0].Default)
	}
}

func TestClone_OverwritesExisting(t *testing.T) {
	cfg := makeConfig(
		config.Target{Name: "staging", Vars: []config.VarDef{varDef("FOO", "new")}},
		config.Target{Name: "production", Vars: []config.VarDef{varDef("FOO", "old")}},
	)
	res, err := Clone(cfg, "staging", "production", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Overwritten != 1 {
		t.Errorf("expected 1 overwritten, got %d", res.Overwritten)
	}
	prod := findTarget(cfg, "production")
	if prod.Vars[0].Default != "new" {
		t.Errorf("expected value 'new', got %q", prod.Vars[0].Default)
	}
}

func TestClone_FiltersByKey(t *testing.T) {
	cfg := makeConfig(
		config.Target{Name: "staging", Vars: []config.VarDef{varDef("FOO", "1"), varDef("BAR", "2")}},
	)
	res, err := Clone(cfg, "staging", "production", Options{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Copied != 1 {
		t.Errorf("expected 1 copied, got %d", res.Copied)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestClone_UnknownSourceErrors(t *testing.T) {
	cfg := makeConfig()
	_, err := Clone(cfg, "ghost", "production", Options{})
	if err == nil {
		t.Error("expected error for unknown source target")
	}
}
