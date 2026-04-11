package mask

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

func TestApply_SensitiveFullMask(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "API_KEY", Default: "supersecret", Sensitive: true},
	})
	results, err := Apply(cfg, "production", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Masked != "********" {
		t.Errorf("expected full mask, got %q", results[0].Masked)
	}
}

func TestApply_NonSensitiveUnmasked(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "APP_ENV", Default: "production", Sensitive: false},
	})
	results, err := Apply(cfg, "production", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Masked != "production" {
		t.Errorf("expected plain value, got %q", results[0].Masked)
	}
}

func TestApply_PartialMaskViaExtra(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "DB_PASS", Default: "abcdefgh", Sensitive: false},
	})
	extras := []Rule{{Key: "DB_PASS", Strategy: "partial"}}
	results, err := Apply(cfg, "production", extras)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Masked != "ab****gh" {
		t.Errorf("expected partial mask, got %q", results[0].Masked)
	}
}

func TestApply_HashStrategy(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "TOKEN", Default: "mytoken", Sensitive: false},
	})
	extras := []Rule{{Key: "TOKEN", Strategy: "hash"}}
	results, err := Apply(cfg, "production", extras)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Masked) == 0 || results[0].Masked == "mytoken" {
		t.Errorf("expected hashed value, got %q", results[0].Masked)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig(nil)
	_, err := Apply(cfg, "staging", nil)
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}
