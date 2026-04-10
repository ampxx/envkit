package template

import (
	"os"
	"path/filepath"
	"testing"

	"envkit/internal/config"
)

func makeConfig() *config.Config {
	return &config.Config{
		Vars: []config.VarDef{
			{Key: "APP_NAME", Default: "myapp"},
			{Key: "PORT", Default: "8080"},
			{Key: "DEBUG", Default: "false"},
		},
		Targets: []config.Target{
			{
				Name: "production",
				Overrides: []config.Override{
					{Key: "PORT", Value: "443"},
					{Key: "DEBUG", Value: "false"},
				},
			},
			{
				Name: "development",
				Overrides: []config.Override{
					{Key: "DEBUG", Value: "true"},
				},
			},
		},
	}
}

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "tpl.tmpl")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp template: %v", err)
	}
	return p
}

func TestGenerateString_BasicSubstitution(t *testing.T) {
	cfg := makeConfig()
	res, err := GenerateString(cfg, "development", "app={{.APP_NAME}} port={{.PORT}} debug={{.DEBUG}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "app=myapp port=8080 debug=true"
	if res.Output != want {
		t.Errorf("got %q, want %q", res.Output, want)
	}
}

func TestGenerateString_TargetOverride(t *testing.T) {
	cfg := makeConfig()
	res, err := GenerateString(cfg, "production", "port={{.PORT}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "port=443" {
		t.Errorf("got %q, want %q", res.Output, "port=443")
	}
}

func TestGenerateString_UnknownTarget(t *testing.T) {
	cfg := makeConfig()
	_, err := GenerateString(cfg, "staging", "{{.APP_NAME}}")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestGenerateString_MissingKey(t *testing.T) {
	cfg := makeConfig()
	_, err := GenerateString(cfg, "development", "{{.UNKNOWN_KEY}}")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestGenerate_FromFile(t *testing.T) {
	cfg := makeConfig()
	p := writeTempTemplate(t, "APP={{.APP_NAME}}\nPORT={{.PORT}}\n")
	res, err := Generate(cfg, "production", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "APP=myapp\nPORT=443\n"
	if res.Output != want {
		t.Errorf("got %q, want %q", res.Output, want)
	}
}

func TestGenerate_FileNotFound(t *testing.T) {
	cfg := makeConfig()
	_, err := Generate(cfg, "development", "/nonexistent/path/tpl.tmpl")
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}
