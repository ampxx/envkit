package envimport_test

import (
	"os"
	"path/filepath"
	"testing"

	"envkit/internal/config"
	"envkit/internal/envimport"
)

func makeConfig(targetName string, vars ...config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: targetName, Vars: vars},
		},
	}
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestApply_ImportsNewKeys(t *testing.T) {
	cfg := makeConfig("staging")
	envFile := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	results, err := envimport.Apply(cfg, envFile, envimport.Options{Target: "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != "imported" {
			t.Errorf("key %q: expected imported, got %s", r.Key, r.Status)
		}
	}
}

func TestApply_SkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := makeConfig("staging", config.VarDef{Key: "FOO", Value: "original"})
	envFile := writeTempEnv(t, "FOO=new\n")

	results, err := envimport.Apply(cfg, envFile, envimport.Options{Target: "staging"})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Status != "skipped" {
		t.Errorf("expected skipped, got %s", results[0].Status)
	}
	if cfg.Targets[0].Vars[0].Value != "original" {
		t.Error("value should not have changed")
	}
}

func TestApply_OverwritesExisting(t *testing.T) {
	cfg := makeConfig("prod", config.VarDef{Key: "FOO", Value: "old"})
	envFile := writeTempEnv(t, "FOO=new\n")

	_, err := envimport.Apply(cfg, envFile, envimport.Options{Target: "prod", Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Targets[0].Vars[0].Value != "new" {
		t.Errorf("expected 'new', got %q", cfg.Targets[0].Vars[0].Value)
	}
}

func TestApply_DryRunDoesNotMutate(t *testing.T) {
	cfg := makeConfig("staging")
	envFile := writeTempEnv(t, "FOO=bar\n")

	_, err := envimport.Apply(cfg, envFile, envimport.Options{Target: "staging", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Targets[0].Vars) != 0 {
		t.Error("dry-run should not mutate config")
	}
}

func TestApply_FiltersByKey(t *testing.T) {
	cfg := makeConfig("staging")
	envFile := writeTempEnv(t, "FOO=1\nBAR=2\nBAZ=3\n")

	results, err := envimport.Apply(cfg, envFile, envimport.Options{
		Target: "staging",
		Keys:   []string{"FOO", "BAZ"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig("staging")
	envFile := writeTempEnv(t, "FOO=bar\n")

	_, err := envimport.Apply(cfg, envFile, envimport.Options{Target: "prod"})
	if err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestApply_MissingEnvFile(t *testing.T) {
	cfg := makeConfig("staging")
	_, err := envimport.Apply(cfg, filepath.Join(t.TempDir(), "missing.env"), envimport.Options{Target: "staging"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
