package unused_test

import (
	"testing"

	"envkit/internal/config"
	"envkit/internal/unused"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func varDef(key, defaultVal string) config.VarDef {
	return config.VarDef{Key: key, Default: defaultVal}
}

func TestCheck_NoIssues(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "production",
		Vars: []config.VarDef{varDef("API_KEY", "")},
	})
	envMaps := map[string]map[string]string{
		"production": {"API_KEY": "secret"},
	}

	report := unused.Check(cfg, envMaps)

	if report.HasIssues() {
		t.Fatalf("expected no issues, got: %+v", report.Results)
	}
}

func TestCheck_MissingFromEnvNoDefault(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "staging",
		Vars: []config.VarDef{varDef("DB_PASS", "")},
	})
	envMaps := map[string]map[string]string{
		"staging": {},
	}

	report := unused.Check(cfg, envMaps)

	if !report.HasIssues() {
		t.Fatal("expected issues but got none")
	}
	if report.Results[0].Key != "DB_PASS" {
		t.Errorf("unexpected key: %s", report.Results[0].Key)
	}
}

func TestCheck_HasDefaultSkipsIssue(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "dev",
		Vars: []config.VarDef{varDef("LOG_LEVEL", "info")},
	})
	envMaps := map[string]map[string]string{
		"dev": {},
	}

	report := unused.Check(cfg, envMaps)

	if report.HasIssues() {
		t.Fatalf("expected no issues for var with default, got: %+v", report.Results)
	}
}

func TestCheck_UndeclaredInEnv(t *testing.T) {
	cfg := makeConfig(config.Target{
		Name: "production",
		Vars: []config.VarDef{},
	})
	envMaps := map[string]map[string]string{
		"production": {"GHOST_VAR": "value"},
	}

	report := unused.Check(cfg, envMaps)

	if !report.HasIssues() {
		t.Fatal("expected issue for undeclared env key")
	}
	if report.Results[0].Key != "GHOST_VAR" {
		t.Errorf("unexpected key: %s", report.Results[0].Key)
	}
}

func TestSummary_NoIssues(t *testing.T) {
	r := unused.Report{}
	if r.Summary() != "no unused or undeclared variables found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestSummary_WithIssues(t *testing.T) {
	r := unused.Report{
		Results: []unused.Result{
			{Target: "prod", Key: "X", Reason: "missing"},
			{Target: "prod", Key: "Y", Reason: "undeclared"},
		},
	}
	if r.Summary() != "2 issue(s) found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}
