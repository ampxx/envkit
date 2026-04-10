package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single environment variable key-value pair.
type Entry struct {
	Key   string
	Value string
}

// ParseFile reads a .env file and returns a slice of entries.
// Lines starting with '#' and empty lines are ignored.
func ParseFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("env: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"`)
		if key == "" {
			continue
		}
		entries = append(entries, Entry{Key: key, Value: value})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %q: %w", path, err)
	}
	return entries, nil
}

// ToMap converts a slice of entries into a map for quick lookup.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// WriteFile writes entries to a .env file, one KEY=VALUE per line.
func WriteFile(path string, entries []Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("env: create %q: %w", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range entries {
		val := e.Value
		if strings.ContainsAny(val, " \t") {
			val = `"` + val + `"`
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return fmt.Errorf("env: write entry: %w", err)
		}
	}
	return w.Flush()
}
