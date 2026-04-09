package schema

import (
	"testing"

	"github.com/envkit/envkit/internal/config"
)

func makeConfig(vars []config.VarDef) *config.Config {
	return &config.Config{
		Targets: map[string]config.Target{
			"dev": {Vars: vars},
		},
	}
}

func TestInfer_URL(t *testing.T) {
	if got := Infer("API_URL", ""); got != TypeURL {
		t.Errorf("expected TypeURL, got %s", got)
	}
}

func TestInfer_Bool(t *testing.T) {
	if got := Infer("ENABLE_FEATURE", ""); got != TypeBool {
		t.Errorf("expected TypeBool, got %s", got)
	}
}

func TestInfer_Int(t *testing.T) {
	if got := Infer("PORT", "8080"); got != TypeInt {
		t.Errorf("expected TypeInt, got %s", got)
	}
}

func TestInfer_Email(t *testing.T) {
	if got := Infer("ADMIN_EMAIL", ""); got != TypeEmail {
		t.Errorf("expected TypeEmail, got %s", got)
	}
}

func TestValidate_NoIssues(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "API_URL", Default: "https://example.com"},
		{Key: "PORT", Default: "3000"},
		{Key: "ENABLE_LOGS", Default: "true"},
	})
	issues := Validate(cfg, "dev")
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestValidate_BadURL(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "API_URL", Default: "not-a-url"},
	})
	issues := Validate(cfg, "dev")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_URL" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestValidate_BadBool(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "ENABLE_CACHE", Default: "yes"},
	})
	issues := Validate(cfg, "dev")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestValidate_UnknownTarget(t *testing.T) {
	cfg := makeConfig([]config.VarDef{})
	issues := Validate(cfg, "prod")
	if len(issues) != 1 || issues[0].Key != "*" {
		t.Errorf("expected wildcard issue for missing target, got %v", issues)
	}
}

func TestValidate_EmptyDefaultSkipped(t *testing.T) {
	cfg := makeConfig([]config.VarDef{
		{Key: "API_URL", Default: ""},
	})
	issues := Validate(cfg, "dev")
	if len(issues) != 0 {
		t.Errorf("empty default should not produce issues, got %v", issues)
	}
}
