package validator

import (
	"testing"
)

func TestValidate_RequiredPresent(t *testing.T) {
	rules := map[string]Rule{
		"APP_ENV": {Required: true},
	}
	env := map[string]string{"APP_ENV": "production"}

	results := Validate(rules, env)
	if len(results) != 1 || !results[0].Passed {
		t.Errorf("expected APP_ENV to pass, got %+v", results)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	rules := map[string]Rule{
		"DATABASE_URL": {Required: true},
	}
	env := map[string]string{}

	results := Validate(rules, env)
	if len(results) != 1 || results[0].Passed {
		t.Errorf("expected DATABASE_URL to fail, got %+v", results)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	rules := map[string]Rule{
		"PORT": {Pattern: `^\d+$`},
	}

	env := map[string]string{"PORT": "8080"}
	results := Validate(rules, env)
	if !results[0].Passed {
		t.Errorf("expected PORT to pass pattern, got %+v", results[0])
	}

	env["PORT"] = "abc"
	results = Validate(rules, env)
	if results[0].Passed {
		t.Errorf("expected PORT to fail pattern, got %+v", results[0])
	}
}

func TestValidate_AllowedValues(t *testing.T) {
	rules := map[string]Rule{
		"LOG_LEVEL": {Allowed: []string{"debug", "info", "warn", "error"}},
	}

	env := map[string]string{"LOG_LEVEL": "info"}
	results := Validate(rules, env)
	if !results[0].Passed {
		t.Errorf("expected LOG_LEVEL to pass, got %+v", results[0])
	}

	env["LOG_LEVEL"] = "verbose"
	results = Validate(rules, env)
	if results[0].Passed {
		t.Errorf("expected LOG_LEVEL to fail, got %+v", results[0])
	}
}

func TestHasFailures(t *testing.T) {
	passing := []Result{{Key: "A", Passed: true}}
	if HasFailures(passing) {
		t.Error("expected no failures")
	}

	mixed := []Result{{Key: "A", Passed: true}, {Key: "B", Passed: false}}
	if !HasFailures(mixed) {
		t.Error("expected failures to be detected")
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	rules := map[string]Rule{
		"OPTIONAL_VAR": {Required: false, Pattern: `^\w+$`},
	}
	env := map[string]string{}

	results := Validate(rules, env)
	if len(results) != 0 {
		t.Errorf("expected no results for optional missing var, got %+v", results)
	}
}
