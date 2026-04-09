package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "envkit.yaml")

	configContent := `version: "1.0"
environment:
  - name: DATABASE_URL
    description: Database connection string
    required: true
    secret: true
  - name: API_KEY
    required: false
    default: "test-key"
targets:
  - name: dev
    description: Development environment
    variables:
      DATABASE_URL: "postgres://localhost/dev"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", config.Version)
	}

	if len(config.Environment) != 2 {
		t.Errorf("Expected 2 environment variables, got %d", len(config.Environment))
	}

	if len(config.Targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(config.Targets))
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "envkit.yaml")

	config := &Config{
		Version: "1.0",
		Environment: []EnvVar{
			{Name: "TEST_VAR", Required: true},
		},
		Targets: []Target{
			{Name: "prod", Description: "Production"},
		},
	}

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestFindTarget(t *testing.T) {
	config := &Config{
		Targets: []Target{
			{Name: "dev"},
			{Name: "prod"},
		},
	}

	target := config.FindTarget("dev")
	if target == nil || target.Name != "dev" {
		t.Error("Failed to find dev target")
	}

	if config.FindTarget("nonexistent") != nil {
		t.Error("Found nonexistent target")
	}
}
