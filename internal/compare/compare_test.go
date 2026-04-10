package compare

import (
	"testing"

	"github.com/user/envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Config {
	return &config.Config{Targets: targets}
}

func makeTarget(name string, kvs ...string) config.Target {
	t := config.Target{Name: name}
	for i := 0; i+1 < len(kvs); i += 2 {
		t.Vars = append(t.Vars, config.VarDef{Key: kvs[i], Default: kvs[i+1]})
	}
	return t
}

func TestTargets_OnlyInA(t *testing.T) {
	cfg := makeConfig(
		makeTarget("prod", "DB_HOST", "db.prod", "API_KEY", "secret"),
		makeTarget("staging", "DB_HOST", "db.staging"),
	)
	res, err := Targets(cfg, "prod", "staging")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.OnlyInA) != 1 || res.OnlyInA[0] != "API_KEY" {
		t.Errorf("expected OnlyInA=[API_KEY], got %v", res.OnlyInA)
	}
}

func TestTargets_OnlyInB(t *testing.T) {
	cfg := makeConfig(
		makeTarget("prod", "DB_HOST", "db.prod"),
		makeTarget("staging", "DB_HOST", "db.staging", "DEBUG", "true"),
	)
	res, err := Targets(cfg, "prod", "staging")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.OnlyInB) != 1 || res.OnlyInB[0] != "DEBUG" {
		t.Errorf("expected OnlyInB=[DEBUG], got %v", res.OnlyInB)
	}
}

func TestTargets_Differing(t *testing.T) {
	cfg := makeConfig(
		makeTarget("prod", "DB_HOST", "db.prod"),
		makeTarget("staging", "DB_HOST", "db.staging"),
	)
	res, err := Targets(cfg, "prod", "staging")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Differing) != 1 || res.Differing[0].Key != "DB_HOST" {
		t.Errorf("expected one differing key DB_HOST, got %v", res.Differing)
	}
}

func TestTargets_Common(t *testing.T) {
	cfg := makeConfig(
		makeTarget("prod", "LOG_LEVEL", "info"),
		makeTarget("staging", "LOG_LEVEL", "info"),
	)
	res, err := Targets(cfg, "prod", "staging")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Common) != 1 || res.Common[0] != "LOG_LEVEL" {
		t.Errorf("expected Common=[LOG_LEVEL], got %v", res.Common)
	}
}

func TestTargets_UnknownTarget(t *testing.T) {
	cfg := makeConfig(makeTarget("prod"))
	_, err := Targets(cfg, "prod", "ghost")
	if err == nil {
		t.Error("expected error for unknown target")
	}
}
