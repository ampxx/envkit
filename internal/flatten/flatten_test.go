package flatten

import (
	"testing"

	"envkit/internal/config"
)

func makeConfig() *config.Document {
	return &config.Document{
		Vars: []config.VarDef{
			{Key: "APP_NAME", Default: "myapp"},
			{Key: "LOG_LEVEL", Default: "info"},
		},
		Targets: []config.Target{
			{
				Name: "prod",
				Vars: []config.VarDef{
					{Key: "LOG_LEVEL", Default: "warn"},
					{Key: "DB_HOST", Default: "db.prod.local"},
				},
			},
		},
	}
}

func findKey(rs []Result, key string) (Result, bool) {
	for _, r := range rs {
		if r.Key == key {
			return r, true
		}
	}
	return Result{}, false
}

func TestApply_GlobalsOnly(t *testing.T) {
	cfg := makeConfig()
	out, err := Apply(cfg, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	r, ok := findKey(out, "LOG_LEVEL")
	if !ok || r.Value != "info" || r.Source != "global" {
		t.Errorf("expected global LOG_LEVEL=info, got %+v", r)
	}
}

func TestApply_TargetOverridesGlobal(t *testing.T) {
	cfg := makeConfig()
	out, err := Apply(cfg, Options{Target: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := findKey(out, "LOG_LEVEL")
	if !ok || r.Value != "warn" || r.Source != "prod" {
		t.Errorf("expected prod LOG_LEVEL=warn, got %+v", r)
	}
}

func TestApply_UnknownTargetErrors(t *testing.T) {
	cfg := makeConfig()
	_, err := Apply(cfg, Options{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_PrefixApplied(t *testing.T) {
	cfg := makeConfig()
	out, err := Apply(cfg, Options{Prefix: "svc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range out {
		if len(r.Key) < 4 || r.Key[:4] != "SVC_" {
			t.Errorf("expected SVC_ prefix, got %q", r.Key)
		}
	}
}

func TestApply_FilterByKeys(t *testing.T) {
	cfg := makeConfig()
	out, err := Apply(cfg, Options{Keys: []string{"APP_NAME"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "APP_NAME" {
		t.Errorf("expected only APP_NAME, got %+v", out)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	cfg := makeConfig()
	out, err := Apply(cfg, Options{Prefix: "svc", Separator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range out {
		if len(r.Key) < 4 || r.Key[:4] != "SVC." {
			t.Errorf("expected SVC. prefix with dot separator, got %q", r.Key)
		}
	}
}
