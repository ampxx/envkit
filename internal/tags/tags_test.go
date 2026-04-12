package tags

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

func varDef(key string, tags ...string) config.VarDef {
	return config.VarDef{Key: key, Tags: tags}
}

func TestApply_AddsTags(t *testing.T) {
	cfg := makeConfig("prod", []config.VarDef{varDef("DB_URL")})
	results, err := Apply(cfg, Options{Target: "prod", Add: []string{"secret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Action != "updated" {
		t.Errorf("expected action=updated, got %q", results[0].Action)
	}
	if len(results[0].Tags) != 1 || results[0].Tags[0] != "secret" {
		t.Errorf("unexpected tags: %v", results[0].Tags)
	}
}

func TestApply_RemovesTags(t *testing.T) {
	cfg := makeConfig("prod", []config.VarDef{varDef("API_KEY", "secret", "rotate")})
	results, err := Apply(cfg, Options{Target: "prod", Remove: []string{"rotate"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Tags) != 1 || results[0].Tags[0] != "secret" {
		t.Errorf("unexpected tags: %v", results[0].Tags)
	}
}

func TestApply_UnchangedWhenNoOp(t *testing.T) {
	cfg := makeConfig("prod", []config.VarDef{varDef("PORT", "infra")})
	results, err := Apply(cfg, Options{Target: "prod", Add: []string{"infra"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Action != "unchanged" {
		t.Errorf("expected unchanged, got %q", results[0].Action)
	}
}

func TestApply_FiltersByKey(t *testing.T) {
	cfg := makeConfig("prod", []config.VarDef{
		varDef("DB_URL"),
		varDef("API_KEY"),
	})
	results, err := Apply(cfg, Options{Target: "prod", Keys: []string{"DB_URL"}, Add: []string{"secret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_URL" {
		t.Errorf("expected DB_URL, got %q", results[0].Key)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	cfg := makeConfig("prod", nil)
	_, err := Apply(cfg, Options{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestList_ReturnsUniqueSortedTags(t *testing.T) {
	cfg := makeConfig("prod", []config.VarDef{
		varDef("A", "secret", "rotate"),
		varDef("B", "secret", "infra"),
	})
	tags, err := List(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"infra", "rotate", "secret"}
	if len(tags) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], tag)
		}
	}
}
