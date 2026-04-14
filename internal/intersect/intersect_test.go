package intersect_test

import (
	"testing"

	"github.com/your-org/envkit/internal/config"
	"github.com/your-org/envkit/internal/intersect"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func makeTarget(name string, keys ...string) config.Target {
	vars := make([]config.VarDef, len(keys))
	for i, k := range keys {
		vars[i] = config.VarDef{Key: k}
	}
	return config.Target{Name: name, Vars: vars}
}

func TestApply_CommonKeys(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", "DB_URL", "API_KEY", "PORT"),
		makeTarget("production", "DB_URL", "API_KEY", "SECRET"),
	)
	res, err := intersect.Apply(cfg, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.CommonKeys) != 2 {
		t.Errorf("expected 2 common keys, got %d", len(res.CommonKeys))
	}
}

func TestApply_OnlyInA(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", "DB_URL", "PORT"),
		makeTarget("production", "DB_URL"),
	)
	res, err := intersect.Apply(cfg, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.OnlyInA) != 1 || res.OnlyInA[0] != "PORT" {
		t.Errorf("expected OnlyInA=[PORT], got %v", res.OnlyInA)
	}
}

func TestApply_OnlyInB(t *testing.T) {
	cfg := makeConfig(
		makeTarget("staging", "DB_URL"),
		makeTarget("production", "DB_URL", "SECRET"),
	)
	res, err := intersect.Apply(cfg, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.OnlyInB) != 1 || res.OnlyInB[0] != "SECRET" {
		t.Errorf("expected OnlyInB=[SECRET], got %v", res.OnlyInB)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig(makeTarget("staging", "DB_URL"))
	_, err := intersect.Apply(cfg, "staging", "missing")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_NoCommon(t *testing.T) {
	cfg := makeConfig(
		makeTarget("a", "FOO"),
		makeTarget("b", "BAR"),
	)
	res, err := intersect.Apply(cfg, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.CommonKeys) != 0 {
		t.Errorf("expected 0 common keys, got %d", len(res.CommonKeys))
	}
}
