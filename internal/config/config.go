package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// EnvVar represents a single environment variable configuration
type EnvVar struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Required    bool     `yaml:"required"`
	Default     string   `yaml:"default,omitempty"`
	Targets     []string `yaml:"targets,omitempty"`
	Secret      bool     `yaml:"secret,omitempty"`
}

// Target represents a deployment target (e.g., dev, staging, prod)
type Target struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Variables   map[string]string `yaml:"variables,omitempty"`
}

// Config represents the complete envkit configuration
type Config struct {
	Version     string   `yaml:"version"`
	Environment []EnvVar `yaml:"environment"`
	Targets     []Target `yaml:"targets"`
}

// Load reads and parses the envkit configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save writes the configuration to a file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// FindTarget returns a target by name
func (c *Config) FindTarget(name string) *Target {
	for i := range c.Targets {
		if c.Targets[i].Name == name {
			return &c.Targets[i]
		}
	}
	return nil
}

// FindEnvVar returns an environment variable by name
func (c *Config) FindEnvVar(name string) *EnvVar {
	for i := range c.Environment {
		if c.Environment[i].Name == name {
			return &c.Environment[i]
		}
	}
	return nil
}
