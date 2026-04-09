package render

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"envkit/internal/config"
)

// Result holds the output of a rendered template.
type Result struct {
	Target string
	Output string
}

// Render applies Go template substitution to a template string using
// the resolved values from the given target in cfg.
func Render(cfg *config.Config, targetName, tmpl string) (*Result, error) {
	target := cfg.FindTarget(targetName)
	if target == nil {
		return nil, fmt.Errorf("target %q not found", targetName)
	}

	vals := buildValues(cfg, target)

	t, err := template.New("envkit").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, vals); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return &Result{Target: targetName, Output: buf.String()}, nil
}

// RenderFile reads a template from disk and renders it.
func RenderFile(cfg *config.Config, targetName, templatePath string) (*Result, error) {
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("read template file: %w", err)
	}
	return Render(cfg, targetName, string(raw))
}

// buildValues constructs the data map passed to the template engine.
// Global vars are included first; target-level overrides take precedence.
func buildValues(cfg *config.Config, target *config.Target) map[string]string {
	vals := make(map[string]string)
	for _, v := range cfg.Vars {
		vals[strings.ToUpper(v.Key)] = v.Default
	}
	for _, v := range target.Overrides {
		vals[strings.ToUpper(v.Key)] = v.Value
	}
	return vals
}
