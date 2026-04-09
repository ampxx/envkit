package render

import (
	"os"
	"path/filepath"
	"testing"

	"envkit/internal/config"
)

func makeConfig() *config.Config {
	return &config.Config{
		Vars: []config.VarDef{
			{Key: "APP_NAME", Default: "envkit"},
			{Key: "LOG_LEVEL", Default: "info"},
		},
		Targets: []config.Target{
			{
				Name: "production",
				Overrides: []config.Override{
					{Key: "LOG_LEVEL", Value: "warn"},
				},
			},
			{
				Name: "development",
				Overrides: []config.Override{},
			},
		},
	}
}

func TestRender_BasicSubstitution(t *testing.T) {
	cfg := makeConfig()
	res, err := Render(cfg, "development", "app={{ .APP_NAME }} level={{ .LOG_LEVEL }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "app=envkit level=info"
	if res.Output != want {
		t.Errorf("got %q, want %q", res.Output, want)
	}
}

func TestRender_TargetOverride(t *testing.T) {
	cfg := makeConfig()
	res, err := Render(cfg, "production", "level={{ .LOG_LEVEL }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "level=warn" {
		t.Errorf("got %q, want %q", res.Output, "level=warn")
	}
}

func TestRender_UnknownTarget(t *testing.T) {
	cfg := makeConfig()
	_, err := Render(cfg, "staging", "{{ .APP_NAME }}")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestRender_MissingKey(t *testing.T) {
	cfg := makeConfig()
	_, err := Render(cfg, "development", "{{ .UNDEFINED_VAR }}")
	if err == nil {
		t.Fatal("expected error for missing key in template")
	}
}

func TestRenderFile_ReadsAndRenders(t *testing.T) {
	cfg := makeConfig()
	dir := t.TempDir()
	path := filepath.Join(dir, "tmpl.txt")
	if err := os.WriteFile(path, []byte("name={{ .APP_NAME }}"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := RenderFile(cfg, "development", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "name=envkit" {
		t.Errorf("got %q", res.Output)
	}
}
