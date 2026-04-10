package resolve_test

import (
	"testing"

	"github.com/your-org/envkit/internal/config"
	"github.com/your-org/envkit/internal/resolve"
)

func makeConfig() *config.Config {
	return &config.Config{
		Vars: []config.VarDef{
			{Name: "APP_ENV", Default: "development"},
			{Name: "LOG_LEVEL", Default: "info"},
			{Name: "SECRET_KEY"},
		},
		Targets: []config.Target{
			{
				Name: "production",
				Vars: []config.VarDef{
					{Name: "APP_ENV", Default: "production"},
					{Name: "LOG_LEVEL", Default: "warn"},
				},
			},
		},
	}
}

func TestResolve_GlobalDefaults(t *testing.T) {
	cfg := makeConfig()
	results, err := resolve.Resolve(cfg, resolve.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := resolve.ToMap(results)
	if m["APP_ENV"] != "development" {
		t.Errorf("expected development, got %s", m["APP_ENV"])
	}
	if m["LOG_LEVEL"] != "info" {
		t.Errorf("expected info, got %s", m["LOG_LEVEL"])
	}
}

func TestResolve_TargetOverridesDefault(t *testing.T) {
	cfg := makeConfig()
	results, err := resolve.Resolve(cfg, resolve.Options{Target: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := resolve.ToMap(results)
	if m["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", m["APP_ENV"])
	}
	if m["LOG_LEVEL"] != "warn" {
		t.Errorf("expected warn, got %s", m["LOG_LEVEL"])
	}
}

func TestResolve_CallerOverrideWins(t *testing.T) {
	cfg := makeConfig()
	results, err := resolve.Resolve(cfg, resolve.Options{
		Target:   "production",
		Override: map[string]string{"LOG_LEVEL": "debug"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := resolve.ToMap(results)
	if m["LOG_LEVEL"] != "debug" {
		t.Errorf("expected debug, got %s", m["LOG_LEVEL"])
	}
}

func TestResolve_UnknownTargetErrors(t *testing.T) {
	cfg := makeConfig()
	_, err := resolve.Resolve(cfg, resolve.Options{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestFormat_ShowSource(t *testing.T) {
	cfg := makeConfig()
	results, _ := resolve.Resolve(cfg, resolve.Options{})
	out := resolve.Format(results, true)
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
	// source annotation should appear
	if !contains(out, "[default]") {
		t.Errorf("expected [default] source annotation in output")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
