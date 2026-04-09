package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompter handles interactive user input for missing environment variables.
type Prompter struct {
	reader *bufio.Reader
}

// New creates a new Prompter reading from stdin.
func New() *Prompter {
	return &Prompter{reader: bufio.NewReader(os.Stdin)}
}

// AskString prompts the user for a string value with an optional default.
func (p *Prompter) AskString(key, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("  %s [%s]: ", key, defaultVal)
	} else {
		fmt.Printf("  %s: ", key)
	}

	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading input for %s: %w", key, err)
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal, nil
	}
	return line, nil
}

// AskConfirm prompts the user for a yes/no confirmation.
func (p *Prompter) AskConfirm(question string) (bool, error) {
	fmt.Printf("  %s [y/N]: ", question)

	line, err := p.reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("reading confirmation: %w", err)
	}

	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes", nil
}

// FillMissing interactively prompts for each missing key and returns a map of filled values.
func (p *Prompter) FillMissing(missingKeys []string) (map[string]string, error) {
	if len(missingKeys) == 0 {
		return nil, nil
	}

	fmt.Println("\nFill in missing environment variables (press Enter to skip):")

	result := make(map[string]string, len(missingKeys))
	for _, key := range missingKeys {
		val, err := p.AskString(key, "")
		if err != nil {
			return nil, err
		}
		if val != "" {
			result[key] = val
		}
	}
	return result, nil
}
