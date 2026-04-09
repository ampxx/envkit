package schema

import (
	"fmt"
	"strings"

	"github.com/envkit/envkit/internal/config"
)

// FieldType represents the expected type of an environment variable value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeURL     FieldType = "url"
	TypeEmail   FieldType = "email"
)

// SchemaIssue describes a single schema violation.
type SchemaIssue struct {
	Key     string
	Message string
}

func (s SchemaIssue) String() string {
	return fmt.Sprintf("[%s] %s", s.Key, s.Message)
}

// Infer derives a FieldType from a variable's name and default value heuristics.
func Infer(key, value string) FieldType {
	lower := strings.ToLower(key)
	switch {
	case strings.HasSuffix(lower, "_url") || strings.HasPrefix(value, "http"):
		return TypeURL
	case strings.HasSuffix(lower, "_email") || strings.Contains(value, "@"):
		return TypeEmail
	case value == "true" || value == "false" ||
		strings.HasPrefix(lower, "enable_") || strings.HasPrefix(lower, "disable_"):
		return TypeBool
	case isNumeric(value):
		return TypeInt
	default:
		return TypeString
	}
}

// Validate checks all variables in a config target against their inferred or
// declared schema types and returns a list of issues.
func Validate(cfg *config.Config, target string) []SchemaIssue {
	var issues []SchemaIssue

	tgt, ok := cfg.Targets[target]
	if !ok {
		return []SchemaIssue{{Key: "*", Message: fmt.Sprintf("target %q not found", target)}}
	}

	for _, v := range tgt.Vars {
		ft := Infer(v.Key, v.Default)
		if err := checkType(v.Key, v.Default, ft); err != nil {
			issues = append(issues, SchemaIssue{Key: v.Key, Message: err.Error()})
		}
	}
	return issues
}

func checkType(key, value string, ft FieldType) error {
	if value == "" {
		return nil
	}
	switch ft {
	case TypeBool:
		if value != "true" && value != "false" {
			return fmt.Errorf("expected bool (true/false), got %q", value)
		}
	case TypeInt:
		if !isNumeric(value) {
			return fmt.Errorf("expected integer, got %q", value)
		}
	case TypeURL:
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return fmt.Errorf("expected URL starting with http:// or https://, got %q", value)
		}
	case TypeEmail:
		if !strings.Contains(value, "@") {
			return fmt.Errorf("expected email address, got %q", value)
		}
	}
	return nil
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
