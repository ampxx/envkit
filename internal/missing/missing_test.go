package missing_test

import (
	"os"
	"path/filepath"
	"testing"

	"envkit/internal/config"
	"envkit/internal/missing"
)

func makeConfig(target string, keys ...string) *config.Document {
	vars := make([]config.VarDef, len(keys))
	for i, k := range keys {
		vars[i] = config.VarDef{Key: k, Required: true}
	}
	return &config.Document{
		Targets: []config.Target{{Name: target, Vars: vars}},
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
	return filepath.Clean(f.Name())
}

func TestCheck_NoMissing(t *testing.T) {
	cfg := makeConfig("prod", "DB_HOST", "DB_PORT")
	envFile := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	rep, err := missing.Check(cfg, "prod", envFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.HasMissing() {
		t.Errorf("expected no missing vars, got %v", rep.Results)
	}
}

func TestCheck_SomeMissing(t *testing.T) {
	cfg := makeConfig("prod", "DB_HOST", "DB_PORT", "API_KEY")
	envFile := writeTempEnv(t, "DB_HOST=localhost\n")

	rep, err := missing.Check(cfg, "prod", envFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 2 {
		t.Fatalf("expected 2 missing, got %d", len(rep.Results))
	}
}

func TestCheck_UnknownTarget(t *testing.T) {
	cfg := makeConfig("prod", "DB_HOST")
	envFile := writeTempEnv(t, "DB_HOST=localhost\n")

	_, err := missing.Check(cfg, "staging", envFile)
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestCheck_RequiredMissing(t *testing.T) {
	cfg := &config.Document{
		Targets: []config.Target{{
			Name: "prod",
			Vars: []config.VarDef{
				{Key: "DB_HOST", Required: true},
				{Key: "LOG_LEVEL", Required: false},
			},
		}},
	}
	envFile := writeTempEnv(t, "")

	rep, err := missing.Check(cfg, "prod", envFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req := rep.RequiredMissing()
	if len(req) != 1 || req[0].Key != "DB_HOST" {
		t.Errorf("expected only DB_HOST in required missing, got %v", req)
	}
}
