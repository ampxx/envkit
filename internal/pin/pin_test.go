package pin

import (
	"testing"
	"time"

	"envkit/internal/config"
)

func makeConfig(targets ...config.Target) *config.Document {
	return &config.Document{Targets: targets}
}

func makeTarget(name string, vars ...config.VarDef) config.Target {
	return config.Target{Name: name, Vars: vars}
}

func varDef(key, def string) config.VarDef {
	return config.VarDef{Key: key, Default: def}
}

func TestPin_AllKeys(t *testing.T) {
	cfg := makeConfig(makeTarget("prod",
		varDef("DB_HOST", "db.prod.example.com"),
		varDef("PORT", "5432"),
	))
	res, err := Pin(cfg, "prod", nil, "ci-bot", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Errorf("expected 2 pinned, got %d", len(res.Pinned))
	}
}

func TestPin_FilterByKey(t *testing.T) {
	cfg := makeConfig(makeTarget("prod",
		varDef("DB_HOST", "db.prod.example.com"),
		varDef("PORT", "5432"),
	))
	res, err := Pin(cfg, "prod", []string{"DB_HOST"}, "ci-bot", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 1 || res.Pinned[0].Key != "DB_HOST" {
		t.Errorf("expected only DB_HOST pinned, got %+v", res.Pinned)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "PORT" {
		t.Errorf("expected PORT skipped, got %+v", res.Skipped)
	}
}

func TestPin_UnknownTarget(t *testing.T) {
	cfg := makeConfig(makeTarget("staging", varDef("X", "1")))
	_, err := Pin(cfg, "prod", nil, "me", 0)
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestPin_NoValueSkipsWithError(t *testing.T) {
	cfg := makeConfig(makeTarget("prod",
		config.VarDef{Key: "EMPTY_VAR"},
	))
	res, err := Pin(cfg, "prod", nil, "me", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Errors) != 1 {
		t.Errorf("expected 1 error entry, got %d", len(res.Errors))
	}
}

func TestPin_WithTTL(t *testing.T) {
	cfg := makeConfig(makeTarget("prod", varDef("KEY", "val")))
	res, err := Pin(cfg, "prod", nil, "me", 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) == 0 {
		t.Fatal("expected pinned entry")
	}
	if res.Pinned[0].ExpiresAt.IsZero() {
		t.Error("expected ExpiresAt to be set")
	}
}

func TestExpired_ReturnsExpiredEntries(t *testing.T) {
	now := time.Now().UTC()
	entries := []Entry{
		{Key: "OLD", ExpiresAt: now.Add(-time.Hour)},
		{Key: "FRESH", ExpiresAt: now.Add(time.Hour)},
		{Key: "NOTTL"},
	}
	exp := Expired(entries)
	if len(exp) != 1 || exp[0].Key != "OLD" {
		t.Errorf("expected only OLD expired, got %+v", exp)
	}
}
