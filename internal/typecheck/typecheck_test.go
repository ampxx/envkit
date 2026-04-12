package typecheck

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(vars []config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: "production", Vars: vars},
		},
	}
}

func varDef(key, typ string) config.VarDef {
	return config.VarDef{Key: key, Type: typ}
}

func TestRun_IntValid(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("PORT", "int")})
	r, err := Run(cfg, "production", map[string]string{"PORT": "8080"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Results[0].Status != StatusOK {
		t.Errorf("expected OK, got %s: %s", r.Results[0].Status, r.Results[0].Reason)
	}
}

func TestRun_IntInvalid(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("PORT", "int")})
	r, err := Run(cfg, "production", map[string]string{"PORT": "abc"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Results[0].Status != StatusFail {
		t.Errorf("expected Fail, got %s", r.Results[0].Status)
	}
}

func TestRun_BoolValid(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("DEBUG", "bool")})
	r, _ := Run(cfg, "production", map[string]string{"DEBUG": "true"})
	if r.Results[0].Status != StatusOK {
		t.Errorf("expected OK, got %s", r.Results[0].Status)
	}
}

func TestRun_URLInvalid(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("API_URL", "url")})
	r, _ := Run(cfg, "production", map[string]string{"API_URL": "not-a-url"})
	if r.Results[0].Status != StatusFail {
		t.Errorf("expected Fail, got %s", r.Results[0].Status)
	}
}

func TestRun_EmailValid(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("ADMIN", "email")})
	r, _ := Run(cfg, "production", map[string]string{"ADMIN": "admin@example.com"})
	if r.Results[0].Status != StatusOK {
		t.Errorf("expected OK, got %s", r.Results[0].Status)
	}
}

func TestRun_SkippedWhenNotSet(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("PORT", "int")})
	r, _ := Run(cfg, "production", map[string]string{})
	if r.Results[0].Status != StatusSkipped {
		t.Errorf("expected Skipped, got %s", r.Results[0].Status)
	}
}

func TestRun_UnknownTarget(t *testing.T) {
	cfg := makeConfig([]config.VarDef{varDef("X", "int")})
	_, err := Run(cfg, "staging", map[string]string{})
	if err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestHasFailures_True(t *testing.T) {
	r := Report{Results: []Result{{Status: StatusFail}}}
	if !HasFailures(r) {
		t.Error("expected HasFailures to return true")
	}
}

func TestHasFailures_False(t *testing.T) {
	r := Report{Results: []Result{{Status: StatusOK}}}
	if HasFailures(r) {
		t.Error("expected HasFailures to return false")
	}
}
