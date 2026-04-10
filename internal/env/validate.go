package env

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("line %d [%s]: %s", e.Line, e.Key, e.Message)
}

// ValidateFile parses an env file and checks for common issues:
// duplicate keys, empty keys, keys with spaces, and values that look
// like unquoted multiline content.
func ValidateFile(path string) ([]ValidationError, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	var errs []ValidationError
	seen := make(map[string]int)

	for i, e := range entries {
		lineNum := i + 1

		if e.Key == "" {
			errs = append(errs, ValidationError{Line: lineNum, Key: "(empty)", Message: "key must not be empty"})
			continue
		}

		if strings.Contains(e.Key, " ") {
			errs = append(errs, ValidationError{Line: lineNum, Key: e.Key, Message: "key must not contain spaces"})
		}

		if strings.ToLower(e.Key) == e.Key && strings.ContainsAny(e.Key, "abcdefghijklmnopqrstuvwxyz") {
			errs = append(errs, ValidationError{Line: lineNum, Key: e.Key, Message: "key should be uppercase"})
		}

		if prev, dup := seen[e.Key]; dup {
			errs = append(errs, ValidationError{
				Line:    lineNum,
				Key:     e.Key,
				Message: fmt.Sprintf("duplicate key (first seen on line %d)", prev),
			})
		} else {
			seen[e.Key] = lineNum
		}

		if strings.Contains(e.Value, "\n") {
			errs = append(errs, ValidationError{Line: lineNum, Key: e.Key, Message: "value contains unquoted newline"})
		}
	}

	return errs, nil
}

// ValidationSummary returns a human-readable summary string.
func ValidationSummary(errs []ValidationError) string {
	if len(errs) == 0 {
		return "no issues found"
	}
	return fmt.Sprintf("%d issue(s) found", len(errs))
}
