package merge

import (
	"testing"

	"github.com/user/envkit/internal/config"
)

func makeConfig(targets map[string][]config.VarDef) *config.Config {
	c := &config.Config{
		Version: "1",
		Targets: make(map[string]config.Target),
	}
	for name, vars := range targets {
		c.Targets[name] = config.Target{Vars: vars}
	}
	return c
}

func varDef(name, def string) config.VarDef {
	return config.VarDef{Name: name, Default: def}
}

func TestMerge_AddsNewTarget(t *testing.T) {
	base := makeConfig(map[string][]config.VarDef{"prod": {varDef("A", "1")}})
	other := makeConfig(map[string][]config.VarDef{"staging": {varDef("B", "2")}})

	res, err := Merge(base, other, StrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Merged.Targets["staging"]; !ok {
		t.Error("expected staging target to be present")
	}
}

func TestMerge_StrategyOurs_KeepsBase(t *testing.T) {
	base := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "8080")}})
	other := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "9090")}})

	res, err := Merge(base, other, StrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	got := res.Merged.Targets["prod"].Vars[0].Default
	if got != "8080" {
		t.Errorf("expected 8080, got %s", got)
	}
}

func TestMerge_StrategyTheirs_TakesOther(t *testing.T) {
	base := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "8080")}})
	other := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "9090")}})

	res, err := Merge(base, other, StrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Merged.Targets["prod"].Vars[0].Default
	if got != "9090" {
		t.Errorf("expected 9090, got %s", got)
	}
}

func TestMerge_StrategyPrompt_ReturnsConflicts(t *testing.T) {
	base := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "8080")}})
	other := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "9090")}})

	res, err := Merge(base, other, StrategyPrompt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Key != "PORT" {
		t.Errorf("expected conflict key PORT, got %s", res.Conflicts[0].Key)
	}
}

func TestMerge_NilConfig_ReturnsError(t *testing.T) {
	_, err := Merge(nil, &config.Config{}, StrategyOurs)
	if err == nil {
		t.Error("expected error for nil base config")
	}
}

func TestApplyResolutions(t *testing.T) {
	base := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "8080")}})
	other := makeConfig(map[string][]config.VarDef{"prod": {varDef("PORT", "9090")}})

	res, _ := Merge(base, other, StrategyPrompt)
	ApplyResolutions(res, map[string]string{"PORT": "7070"})

	got := res.Merged.Targets["prod"].Vars[0].Default
	if got != "7070" {
		t.Errorf("expected 7070, got %s", got)
	}
}
