package trim

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig(targetName string, vars []config.VarDef) *config.Document {
	return &config.Document{
		Targets: []config.Target{
			{Name: targetName, Vars: vars},
		},
	}
}

func varDef(key, def string) config.VarDef {
	return config.VarDef{Key: key, Default: def}
}

func TestApply_TrimSpace(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("FOO", "  hello  "),
		varDef("BAR", "world"),
	})
	results, err := Apply(cfg, "dev", Options{TrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Changed || results[0].Trimmed != "hello" {
		t.Errorf("expected FOO trimmed to 'hello', got %q", results[0].Trimmed)
	}
	if results[1].Changed {
		t.Errorf("expected BAR unchanged")
	}
}

func TestApply_TrimQuotes(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("A", `"quoted"`),
		varDef("B", `'single'`),
		varDef("C", "plain"),
	})
	results, err := Apply(cfg, "dev", Options{TrimQuotes: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Trimmed != "quoted" {
		t.Errorf("expected 'quoted', got %q", results[0].Trimmed)
	}
	if results[1].Trimmed != "single" {
		t.Errorf("expected 'single', got %q", results[1].Trimmed)
	}
	if results[2].Changed {
		t.Errorf("expected C unchanged")
	}
}

func TestApply_FilterByKey(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("FOO", "  a  "),
		varDef("BAR", "  b  "),
	})
	results, err := Apply(cfg, "dev", Options{TrimSpace: true, Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "FOO" {
		t.Errorf("expected FOO, got %q", results[0].Key)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig("dev", nil)
	_, err := Apply(cfg, "prod", Options{})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_BothOptions(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("KEY", `  "padded"  `),
	})
	results, err := Apply(cfg, "dev", Options{TrimSpace: true, TrimQuotes: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Trimmed != "padded" {
		t.Errorf("expected 'padded', got %q", results[0].Trimmed)
	}
}
