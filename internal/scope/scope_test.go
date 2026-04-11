package scope

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func makeTarget(name string, vars ...config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func varDef(key, value, def string) config.VarDef {
	return config.VarDef{Key: key, Value: value, Default: def}
}

func TestApply_ReturnsAllVars(t *testing.T) {
	cfg := makeConfig(makeTarget("prod",
		varDef("DB_HOST", "db.prod", ""),
		varDef("DB_PORT", "", "5432"),
	))
	res, err := Apply(cfg, Options{Target: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["DB_HOST"] != "db.prod" {
		t.Errorf("expected db.prod, got %q", res.Vars["DB_HOST"])
	}
	if res.Vars["DB_PORT"] != "5432" {
		t.Errorf("expected 5432, got %q", res.Vars["DB_PORT"])
	}
}

func TestApply_FiltersByKeys(t *testing.T) {
	cfg := makeConfig(makeTarget("prod",
		varDef("A", "1", ""),
		varDef("B", "2", ""),
		varDef("C", "3", ""),
	))
	res, err := Apply(cfg, Options{Target: "prod", Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["B"]; ok {
		t.Error("B should have been filtered out")
	}
	if res.Vars["A"] != "1" || res.Vars["C"] != "3" {
		t.Error("A and C should be present")
	}
}

func TestApply_ExcludesKeys(t *testing.T) {
	cfg := makeConfig(makeTarget("prod",
		varDef("SECRET", "s3cr3t", ""),
		varDef("HOST", "localhost", ""),
	))
	res, err := Apply(cfg, Options{Target: "prod", Exclude: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["SECRET"]; ok {
		t.Error("SECRET should have been excluded")
	}
	if res.Vars["HOST"] != "localhost" {
		t.Error("HOST should be present")
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig(makeTarget("prod"))
	_, err := Apply(cfg, Options{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestPrefix_PrependsPrefixToKeys(t *testing.T) {
	r := Result{Target: "prod", Vars: map[string]string{"HOST": "localhost", "PORT": "8080"}}
	prefixed := Prefix(r, "app")
	if _, ok := prefixed.Vars["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := prefixed.Vars["APP_PORT"]; !ok {
		t.Error("expected APP_PORT")
	}
	if _, ok := prefixed.Vars["HOST"]; ok {
		t.Error("original HOST should not exist")
	}
}
