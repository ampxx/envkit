package obsolete_test

import (
	"testing"

	"github.com/envkit/envkit/internal/config"
	"github.com/envkit/envkit/internal/obsolete"
)

func makeConfig(globalVars []config.VarDef, targets []config.Target) *config.Document {
	return &config.Document{
		Vars:    globalVars,
		Targets: targets,
	}
}

func varDef(key, def string) config.VarDef {
	return config.VarDef{Key: key, Default: def}
}

func makeTarget(name string, vars []config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func TestCheck_NoObsolete(t *testing.T) {
	cfg := makeConfig(
		[]config.VarDef{varDef("LOG_LEVEL", "info")},
		[]config.Target{makeTarget("production", []config.VarDef{varDef("LOG_LEVEL", "warn")})},
	)
	report, err := obsolete.Check(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasIssues() {
		t.Errorf("expected no issues, got %d", len(report.Results))
	}
}

func TestCheck_DetectsRedundantVar(t *testing.T) {
	cfg := makeConfig(
		[]config.VarDef{varDef("LOG_LEVEL", "info")},
		[]config.Target{makeTarget("staging", []config.VarDef{varDef("LOG_LEVEL", "info")})},
	)
	report, err := obsolete.Check(cfg, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasIssues() {
		t.Fatal("expected issues, got none")
	}
	if report.Results[0].Key != "LOG_LEVEL" {
		t.Errorf("expected key LOG_LEVEL, got %s", report.Results[0].Key)
	}
}

func TestCheck_UnknownTarget(t *testing.T) {
	cfg := makeConfig(nil, nil)
	_, err := obsolete.Check(cfg, "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestCheck_VarNotInGlobals(t *testing.T) {
	cfg := makeConfig(
		[]config.VarDef{},
		[]config.Target{makeTarget("dev", []config.VarDef{varDef("ONLY_IN_TARGET", "value")})},
	)
	report, err := obsolete.Check(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasIssues() {
		t.Errorf("expected no issues, got %d", len(report.Results))
	}
}
