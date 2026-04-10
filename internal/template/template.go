package template

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"envkit/internal/config"
)

// Result holds the output of a template generation.
type Result struct {
	Target string
	Output string
}

// Generate applies a Go text/template file against the resolved env vars
// for a given target and returns the rendered output.
func Generate(cfg *config.Config, targetName, tmplPath string) (*Result, error) {
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}

	values, err := buildValues(cfg, targetName)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("envkit").Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, values); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return &Result{Target: targetName, Output: sb.String()}, nil
}

// GenerateString applies a raw template string instead of a file path.
func GenerateString(cfg *config.Config, targetName, tmplStr string) (*Result, error) {
	values, err := buildValues(cfg, targetName)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("envkit-inline").Option("missingkey=error").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, values); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return &Result{Target: targetName, Output: sb.String()}, nil
}

// buildValues merges base defaults with target-specific overrides into a map.
func buildValues(cfg *config.Config, targetName string) (map[string]string, error) {
	target := cfg.FindTarget(targetName)
	if target == nil {
		return nil, fmt.Errorf("target %q not found", targetName)
	}

	values := make(map[string]string)
	for _, v := range cfg.Vars {
		values[v.Key] = v.Default
	}
	for _, v := range target.Overrides {
		values[v.Key] = v.Value
	}
	return values, nil
}
