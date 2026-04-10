package init

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Template represents a starter envkit config template.
type Template struct {
	Name        string
	Description string
	Targets     []string
	Vars        []VarDef
}

// VarDef is a starter variable definition.
type VarDef struct {
	Key         string
	Description string
	Required    bool
	Example     string
}

// Options controls the scaffold behaviour.
type Options struct {
	Dir        string
	Template   string
	Force      bool
}

var builtinTemplates = map[string]Template{
	"web": {
		Name:        "web",
		Description: "Basic web application",
		Targets:     []string{"development", "staging", "production"},
		Vars: []VarDef{
			{Key: "PORT", Description: "HTTP listen port", Required: true, Example: "8080"},
			{Key: "DATABASE_URL", Description: "Database connection string", Required: true, Example: "postgres://localhost/mydb"},
			{Key: "SECRET_KEY", Description: "Application secret key", Required: true, Example: "changeme"},
			{Key: "LOG_LEVEL", Description: "Logging verbosity", Required: false, Example: "info"},
		},
	},
	"worker": {
		Name:        "worker",
		Description: "Background worker service",
		Targets:     []string{"development", "production"},
		Vars: []VarDef{
			{Key: "QUEUE_URL", Description: "Message queue URL", Required: true, Example: "amqp://localhost"},
			{Key: "CONCURRENCY", Description: "Worker concurrency", Required: false, Example: "4"},
			{Key: "LOG_LEVEL", Description: "Logging verbosity", Required: false, Example: "info"},
		},
	},
}

// ListTemplates returns the names of all available built-in templates.
func ListTemplates() []string {
	names := make([]string, 0, len(builtinTemplates))
	for k := range builtinTemplates {
		names = append(names, k)
	}
	return names
}

// Scaffold creates an envkit.yaml config file in dir using the named template.
func Scaffold(opts Options) error {
	tmpl, ok := builtinTemplates[opts.Template]
	if !ok {
		return fmt.Errorf("unknown template %q; available: %v", opts.Template, ListTemplates())
	}

	dest := filepath.Join(opts.Dir, "envkit.yaml")
	if !opts.Force {
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("%s already exists; use --force to overwrite", dest)
		}
	}

	doc := buildDocument(tmpl)
	out, err := yaml.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return os.WriteFile(dest, out, 0o644)
}

func buildDocument(tmpl Template) map[string]interface{} {
	vars := make([]map[string]interface{}, 0, len(tmpl.Vars))
	for _, v := range tmpl.Vars {
		entry := map[string]interface{}{
			"key":         v.Key,
			"description": v.Description,
			"required":    v.Required,
			"example":     v.Example,
		}
		vars = append(vars, entry)
	}

	targets := make([]map[string]interface{}, 0, len(tmpl.Targets))
	for _, t := range tmpl.Targets {
		targets = append(targets, map[string]interface{}{"name": t, "vars": []interface{}{}})
	}

	return map[string]interface{}{
		"version": 1,
		"vars":    vars,
		"targets": targets,
	}
}
