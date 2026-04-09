package export

import (
	"fmt"
	"os"

	"github.com/envkit/envkit/internal/config"
)

// RunExport is the CLI entry-point for the export command.
// args: [target] [--format dotenv|shell|json] [--out filepath]
func RunExport(cfgPath, target, format, outFile string) error {
	if target == "" {
		return fmt.Errorf("export: target name is required")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("export: failed to load config: %w", err)
	}

	env, err := loadEnvFile(".env")
	if err != nil {
		// non-fatal — proceed with an empty map so callers can still
		// test format output without a real .env file present.
		env = map[string]string{}
	}

	fmt_ := Format(format)
	if fmt_ == "" {
		fmt_ = FormatDotenv
	}

	opts := Options{
		Format:  fmt_,
		Target:  target,
		OutFile: outFile,
	}

	if err := Export(cfg, env, opts); err != nil {
		return err
	}

	if outFile != "" {
		fmt.Fprintf(os.Stderr, "exported %s config to %s\n", target, outFile)
	}
	return nil
}

// loadEnvFile reads a simple KEY=VALUE dotenv file into a map.
func loadEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, line := range splitLines(string(data)) {
		if line == "" || line[0] == '#' {
			continue
		}
		for i, ch := range line {
			if ch == '=' {
				result[line[:i]] = line[i+1:]
				break
			}
		}
	}
	return result, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, ch := range s {
		if ch == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
