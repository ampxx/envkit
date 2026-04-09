package export

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/envkit/envkit/internal/config"
)

// Format represents an output format for exported env vars.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// Options configures the export behaviour.
type Options struct {
	Format  Format
	Target  string
	OutFile string // empty means stdout
}

// Export writes environment variables for the given target to the configured output.
func Export(cfg *config.Config, env map[string]string, opts Options) error {
	target, err := config.FindTarget(cfg, opts.Target)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	var sb strings.Builder

	switch opts.Format {
	case FormatDotenv:
		writeDotenv(&sb, target, env)
	case FormatShell:
		writeShell(&sb, target, env)
	case FormatJSON:
		writeJSON(&sb, target, env)
	default:
		return fmt.Errorf("export: unknown format %q", opts.Format)
	}

	if opts.OutFile == "" {
		fmt.Print(sb.String())
		return nil
	}

	return os.WriteFile(opts.OutFile, []byte(sb.String()), 0o644)
}

func sortedKeys(target *config.Target, env map[string]string) []string {
	seen := map[string]bool{}
	var keys []string
	for _, v := range target.Vars {
		if _, ok := env[v.Key]; ok {
			keys = append(keys, v.Key)
			seen[v.Key] = true
		}
	}
	for k := range env {
		if !seen[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func writeDotenv(sb *strings.Builder, target *config.Target, env map[string]string) {
	for _, k := range sortedKeys(target, env) {
		fmt.Fprintf(sb, "%s=%s\n", k, env[k])
	}
}

func writeShell(sb *strings.Builder, target *config.Target, env map[string]string) {
	for _, k := range sortedKeys(target, env) {
		fmt.Fprintf(sb, "export %s=%q\n", k, env[k])
	}
}

func writeJSON(sb *strings.Builder, target *config.Target, env map[string]string) {
	keys := sortedKeys(target, env)
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(sb, "  %q: %q%s\n", k, env[k], comma)
	}
	sb.WriteString("}\n")
}
