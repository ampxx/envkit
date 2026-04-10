package rotate

import (
	"fmt"
	"strings"
	"testing"

	"envkit/internal/config"
)

func makeConfig(targets []config.Target) *config.Config {
	return &config.Config{Targets: targets}
}

func varDef(key, def string) config.VarDef {
	return config.VarDef{Key: key, Default: def}
}

func suffixFn(suffix string) RotateFn {
	return func(key, old string) (string, error) {
		return old + suffix, nil
	}
}

func TestRotate_AllKeys(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{Name: "prod", Vars: []config.VarDef{varDef("API_KEY", "old"), varDef("DB_PASS", "pass")}},
	})

	results, err := Rotate(cfg, nil, Options{}, suffixFn("_new"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Rotated {
			t.Errorf("key %q should be rotated", r.Key)
		}
	}
	// Config should be mutated.
	if cfg.Targets[0].Vars[0].Default != "old_new" {
		t.Errorf("expected 'old_new', got %q", cfg.Targets[0].Vars[0].Default)
	}
}

func TestRotate_FilterByKey(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{Name: "prod", Vars: []config.VarDef{varDef("API_KEY", "a"), varDef("DB_PASS", "b")}},
	})

	results, err := Rotate(cfg, nil, Options{Keys: []string{"API_KEY"}}, suffixFn("_r"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "API_KEY" {
		t.Errorf("expected only API_KEY in results")
	}
}

func TestRotate_DryRun_DoesNotMutate(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{Name: "staging", Vars: []config.VarDef{varDef("TOKEN", "original")}},
	})

	_, err := Rotate(cfg, nil, Options{DryRun: true}, suffixFn("_changed"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets[0].Vars[0].Default != "original" {
		t.Errorf("dry-run must not mutate config")
	}
}

func TestRotate_RedactHidesValues(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{Name: "prod", Vars: []config.VarDef{varDef("SECRET", "abc")}},
	})

	results, err := Rotate(cfg, nil, Options{Redact: true}, suffixFn("xyz"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if strings.Contains(results[0].OldVal, "abc") {
		t.Errorf("old value should be redacted")
	}
	if strings.Contains(results[0].NewVal, "abcxyz") {
		t.Errorf("new value should be redacted")
	}
}

func TestRotate_FnError_Propagates(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{Name: "prod", Vars: []config.VarDef{varDef("KEY", "v")}},
	})

	_, err := Rotate(cfg, nil, Options{}, func(key, old string) (string, error) {
		return "", fmt.Errorf("rotation service unavailable")
	})
	if err == nil {
		t.Fatal("expected error from RotateFn")
	}
}

func TestRotate_FilterByTarget(t *testing.T) {
	cfg := makeConfig([]config.Target{
		{Name: "prod", Vars: []config.VarDef{varDef("K", "prod_val")}},
		{Name: "dev", Vars: []config.VarDef{varDef("K", "dev_val")}},
	})

	results, err := Rotate(cfg, []string{"prod"}, Options{}, suffixFn("_r"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Target != "prod" {
		t.Errorf("expected only prod target in results")
	}
	// dev target must be untouched.
	if cfg.Targets[1].Vars[0].Default != "dev_val" {
		t.Errorf("dev target should not be mutated")
	}
}
