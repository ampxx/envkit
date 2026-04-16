package envset_test

import (
	"bytes"
	"strings"
	"testing"

	"envkit/internal/envset"
)

func captureReport(results []envset.Result) string {
	var buf bytes.Buffer
	// call unexported helper via the exported one but redirect — use Summary instead
	_ = buf
	return envset.Summary(results)
}

func TestSummary_AllNew(t *testing.T) {
	results := []envset.Result{
		{Key: "A", Created: true},
		{Key: "B", Created: true},
	}
	s := envset.Summary(results)
	if !strings.Contains(s, "2 set") {
		t.Errorf("expected '2 set' in %q", s)
	}
}

func TestSummary_Updated(t *testing.T) {
	results := []envset.Result{
		{Key: "A", Created: false, Skipped: false},
	}
	s := envset.Summary(results)
	if !strings.Contains(s, "1 updated") {
		t.Errorf("expected '1 updated' in %q", s)
	}
}

func TestSummary_Skipped(t *testing.T) {
	results := []envset.Result{
		{Key: "A", Skipped: true},
		{Key: "B", Skipped: true},
	}
	s := envset.Summary(results)
	if !strings.Contains(s, "2 skipped") {
		t.Errorf("expected '2 skipped' in %q", s)
	}
}

func TestSummary_Mixed(t *testing.T) {
	results := []envset.Result{
		{Key: "A", Created: true},
		{Key: "B", Skipped: true},
		{Key: "C", Created: false, Skipped: false},
	}
	s := envset.Summary(results)
	if !strings.Contains(s, "1 set") || !strings.Contains(s, "1 updated") || !strings.Contains(s, "1 skipped") {
		t.Errorf("unexpected summary: %q", s)
	}
}
