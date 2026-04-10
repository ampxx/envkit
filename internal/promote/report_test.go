package promote

import (
	"strings"
	"testing"
)

func TestPrintReport_ShowsPromoted(t *testing.T) {
	results := []Result{
		{Key: "API_URL", From: "staging", To: "production", NewValue: "https://staging.example.com"},
	}
	var buf strings.Builder
	printReportTo(&buf, results, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "ADD") {
		t.Errorf("expected ADD in output, got:\n%s", out)
	}
	if !strings.Contains(out, "API_URL") {
		t.Errorf("expected API_URL in output, got:\n%s", out)
	}
}

func TestPrintReport_ShowsSkipped(t *testing.T) {
	results := []Result{
		{Key: "SECRET", From: "staging", To: "production", Skipped: true, Reason: "already exists"},
	}
	var buf strings.Builder
	printReportTo(&buf, results, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "SKIP") {
		t.Errorf("expected SKIP in output, got:\n%s", out)
	}
}

func TestPrintReport_ShowsUpdate(t *testing.T) {
	results := []Result{
		{Key: "API_URL", From: "staging", To: "production",
			OldValue: "https://old.example.com", NewValue: "https://new.example.com"},
	}
	var buf strings.Builder
	printReportTo(&buf, results, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "UPDATE") {
		t.Errorf("expected UPDATE in output, got:\n%s", out)
	}
}

func TestSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A"},
		{Key: "B", Skipped: true},
		{Key: "C"},
	}
	s := Summary(results)
	if s != "2 promoted, 1 skipped" {
		t.Errorf("unexpected summary: %s", s)
	}
}
