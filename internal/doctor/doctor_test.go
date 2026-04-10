package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Config {
	return &config.Config{Targets: targets}
}

func makeTarget(name string, vars ...config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0644)
	return f.Name()
}

func TestRun_AllPassing(t *testing.T) {
	envFile := writeTempEnv(t, "API_KEY=abc123\n")
	cfg := makeConfig(makeTarget("production",
		config.VarDef{Key: "API_KEY", Required: true},
	))
	report, err := Run(cfg, "production", envFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasFailures() {
		for _, r := range report.Results {
			if !r.Passed {
				t.Errorf("check %q failed: %s", r.Name, r.Message)
			}
		}
	}
}

func TestRun_MissingEnvFile(t *testing.T) {
	cfg := makeConfig(makeTarget("staging"))
	report, err := Run(cfg, "staging", filepath.Join(t.TempDir(), "nonexistent.env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range report.Results {
		if r.Name == "env file exists" && !r.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected env file exists check to fail")
	}
}

func TestRun_DuplicateVar(t *testing.T) {
	envFile := writeTempEnv(t, "DB_URL=postgres://localhost\n")
	cfg := makeConfig(makeTarget("dev",
		config.VarDef{Key: "DB_URL"},
		config.VarDef{Key: "DB_URL"},
	))
	report, err := Run(cfg, "dev", envFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range report.Results {
		if r.Name == "no duplicate vars" && !r.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate vars check to fail")
	}
}

func TestRun_UnknownTarget(t *testing.T) {
	cfg := makeConfig(makeTarget("production"))
	_, err := Run(cfg, "ghost", ".env")
	if err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestHasFailures_False(t *testing.T) {
	r := &Report{
		Results: []CheckResult{
			{Passed: true},
			{Passed: true},
		},
	}
	if r.HasFailures() {
		t.Error("expected no failures")
	}
}
