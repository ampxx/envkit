package graph

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(vars map[string]config.VarDef) *config.Config {
	return &config.Config{
		Vars:    vars,
		Targets: []config.Target{},
	}
}

func TestBuild_NoDeps(t *testing.T) {
	cfg := makeConfig(map[string]config.VarDef{
		"HOST": {Default: "localhost"},
		"PORT": {Default: "8080"},
	})
	g, err := Build(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(g.Nodes))
	}
}

func TestBuild_WithDeps(t *testing.T) {
	cfg := makeConfig(map[string]config.VarDef{
		"BASE_URL": {Default: "http://${HOST}:${PORT}"},
		"HOST":     {Default: "localhost"},
		"PORT":     {Default: "8080"},
	})
	g, err := Build(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps := g.Nodes["BASE_URL"].Deps
	if len(deps) != 2 {
		t.Errorf("expected 2 deps for BASE_URL, got %d", len(deps))
	}
}

func TestOrder_TopologicalSort(t *testing.T) {
	cfg := makeConfig(map[string]config.VarDef{
		"BASE_URL": {Default: "http://${HOST}:${PORT}"},
		"HOST":     {Default: "localhost"},
		"PORT":     {Default: "8080"},
	})
	g, _ := Build(cfg, "")
	order, err := g.Order()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	index := func(k string) int {
		for i, v := range order {
			if v == k {
				return i
			}
		}
		return -1
	}
	if index("BASE_URL") <= index("HOST") || index("BASE_URL") <= index("PORT") {
		t.Errorf("BASE_URL should come after HOST and PORT, order: %v", order)
	}
}

func TestBuild_CycleDetected(t *testing.T) {
	cfg := makeConfig(map[string]config.VarDef{
		"A": {Default: "${B}"},
		"B": {Default: "${A}"},
	})
	_, err := Build(cfg, "")
	if err == nil {
		t.Error("expected cycle error, got nil")
	}
}

func TestExtractRefs(t *testing.T) {
	refs := extractRefs("http://${HOST}:${PORT}/path")
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
	if refs[0] != "HOST" || refs[1] != "PORT" {
		t.Errorf("unexpected refs: %v", refs)
	}
}

func TestExtractRefs_NoRefs(t *testing.T) {
	refs := extractRefs("plainvalue")
	if len(refs) != 0 {
		t.Errorf("expected no refs, got %v", refs)
	}
}
