package patch

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targetName string, vars ...config.VarDef) *config.Config {
	return &config.Config{
		Targets: []config.Target{
			{Name: targetName, Vars: vars},
		},
	}
}

func varDef(key, val string) config.VarDef {
	return config.VarDef{Key: key, Default: val}
}

func TestApply_SetNewKey(t *testing.T) {
	cfg := makeConfig("prod", varDef("A", "1"))
	res, err := Apply(cfg, "prod", []Change{{Op: OpSet, Key: "B", Value: "2"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res[0].Applied {
		t.Errorf("expected applied")
	}
	if cfg.Targets[0].Vars[1].Default != "2" {
		t.Errorf("expected B=2")
	}
}

func TestApply_SetUpdatesExisting(t *testing.T) {
	cfg := makeConfig("prod", varDef("A", "old"))
	_, err := Apply(cfg, "prod", []Change{{Op: OpSet, Key: "A", Value: "new"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets[0].Vars[0].Default != "new" {
		t.Errorf("expected A=new, got %q", cfg.Targets[0].Vars[0].Default)
	}
}

func TestApply_UnsetExistingKey(t *testing.T) {
	cfg := makeConfig("prod", varDef("A", "1"), varDef("B", "2"))
	res, err := Apply(cfg, "prod", []Change{{Op: OpUnset, Key: "A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res[0].Applied {
		t.Errorf("expected applied")
	}
	if len(cfg.Targets[0].Vars) != 1 || cfg.Targets[0].Vars[0].Key != "B" {
		t.Errorf("expected only B to remain")
	}
}

func TestApply_UnsetMissingKey(t *testing.T) {
	cfg := makeConfig("prod", varDef("A", "1"))
	res, _ := Apply(cfg, "prod", []Change{{Op: OpUnset, Key: "Z"}})
	if res[0].Applied {
		t.Errorf("expected not applied for missing key")
	}
}

func TestApply_RenameKey(t *testing.T) {
	cfg := makeConfig("prod", varDef("OLD", "v"))
	res, err := Apply(cfg, "prod", []Change{{Op: OpRename, Key: "OLD", NewKey: "NEW"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res[0].Applied {
		t.Errorf("expected applied")
	}
	if cfg.Targets[0].Vars[0].Key != "NEW" {
		t.Errorf("expected key=NEW, got %q", cfg.Targets[0].Vars[0].Key)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig("prod")
	_, err := Apply(cfg, "staging", []Change{{Op: OpSet, Key: "X", Value: "1"}})
	if err == nil {
		t.Error("expected error for unknown target")
	}
}
