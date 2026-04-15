package envdiff

import (
	"testing"

	"github.com/your-org/envkit/internal/config"
)

func makeConfig(globalVars []config.VarDef, targets []config.Target) *config.Document {
	return &config.Document{Vars: globalVars, Targets: targets}
}

func makeTarget(name string, vars []config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func varDef(key, def string) config.VarDef {
	return config.VarDef{Key: key, Default: def}
}

func TestCompare_Identical(t *testing.T) {
	doc := makeConfig(nil, []config.Target{
		makeTarget("staging", []config.VarDef{varDef("PORT", "8080")}),
		makeTarget("prod", []config.VarDef{varDef("PORT", "8080")}),
	})
	rep, err := Compare(doc, "staging", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Kind != KindIdentical {
		t.Errorf("expected identical result, got %+v", rep.Results)
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	doc := makeConfig(nil, []config.Target{
		makeTarget("a", []config.VarDef{varDef("ONLY_A", "1")}),
		makeTarget("b", nil),
	})
	rep, err := Compare(doc, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Kind != KindOnlyInA {
		t.Errorf("expected only_in_a, got %+v", rep.Results)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	doc := makeConfig(nil, []config.Target{
		makeTarget("a", nil),
		makeTarget("b", []config.VarDef{varDef("ONLY_B", "2")}),
	})
	rep, err := Compare(doc, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Kind != KindOnlyInB {
		t.Errorf("expected only_in_b, got %+v", rep.Results)
	}
}

func TestCompare_Differing(t *testing.T) {
	doc := makeConfig(nil, []config.Target{
		makeTarget("a", []config.VarDef{varDef("HOST", "localhost")}),
		makeTarget("b", []config.VarDef{varDef("HOST", "prod.example.com")}),
	})
	rep, err := Compare(doc, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Kind != KindDiffering {
		t.Errorf("expected differing, got %+v", rep.Results)
	}
	if rep.Results[0].ValueA != "localhost" || rep.Results[0].ValueB != "prod.example.com" {
		t.Errorf("unexpected values: %+v", rep.Results[0])
	}
}

func TestCompare_UnknownTarget(t *testing.T) {
	doc := makeConfig(nil, []config.Target{
		makeTarget("a", nil),
	})
	_, err := Compare(doc, "a", "missing")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestCompare_GlobalDefaultsUsed(t *testing.T) {
	doc := makeConfig(
		[]config.VarDef{varDef("GLOBAL", "shared")},
		[]config.Target{
			makeTarget("a", nil),
			makeTarget("b", nil),
		},
	)
	rep, err := Compare(doc, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Kind != KindIdentical {
		t.Errorf("expected global default to appear as identical, got %+v", rep.Results)
	}
}
