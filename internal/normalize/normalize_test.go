package normalize

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
		varDef("HOST", "  localhost  "),
		varDef("PORT", "8080"),
	})
	rep, err := Apply(cfg, "dev", Rule{TrimSpace: true}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(rep.Results))
	}
	if rep.Results[0].NewValue != "localhost" {
		t.Errorf("expected trimmed value, got %q", rep.Results[0].NewValue)
	}
	if !rep.Results[0].Changed {
		t.Error("expected Changed=true for HOST")
	}
	if rep.Results[1].Changed {
		t.Error("expected Changed=false for PORT")
	}
}

func TestApply_Uppercase(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("ENV", "production"),
	})
	rep, err := Apply(cfg, "dev", Rule{Uppercase: true}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.Results[0].NewValue != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", rep.Results[0].NewValue)
	}
}

func TestApply_TrimQuotes(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("TOKEN", `"abc123"`),
	})
	rep, err := Apply(cfg, "dev", Rule{TrimQuotes: true}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.Results[0].NewValue != "abc123" {
		t.Errorf("expected abc123, got %q", rep.Results[0].NewValue)
	}
}

func TestApply_FilterByKey(t *testing.T) {
	cfg := makeConfig("dev", []config.VarDef{
		varDef("HOST", "  localhost  "),
		varDef("PORT", "  8080  "),
	})
	rep, err := Apply(cfg, "dev", Rule{TrimSpace: true}, []string{"HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rep.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(rep.Results))
	}
	if rep.Results[0].Key != "HOST" {
		t.Errorf("expected HOST, got %q", rep.Results[0].Key)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig("dev", nil)
	_, err := Apply(cfg, "prod", Rule{}, nil)
	if err == nil {
		t.Error("expected error for unknown target")
	}
}
