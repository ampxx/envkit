package inherit

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintReport_ShowsInherited(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", FromTarget: "staging", ToTarget: "production", Value: "db", Skipped: false},
	}
	var buf bytes.Buffer
	printReportTo(&buf, results, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
}

func TestPrintReport_ShowsSkipped(t *testing.T) {
	results := []Result{
		{Key: "PORT", Skipped: true, Reason: "already exists"},
	}
	var buf bytes.Buffer
	printReportTo(&buf, results, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "SKIP") {
		t.Errorf("expected SKIP in output, got: %s", out)
	}
	if !strings.Contains(out, "already exists") {
		t.Errorf("expected reason in output")
	}
}

func TestSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Skipped: false},
		{Key: "B", Skipped: false},
		{Key: "C", Skipped: true},
	}
	s := Summary(results)
	if !strings.Contains(s, "2 inherited") {
		t.Errorf("expected '2 inherited' in %q", s)
	}
	if !strings.Contains(s, "1 skipped") {
		t.Errorf("expected '1 skipped' in %q", s)
	}
}

func TestPrintReport_Header(t *testing.T) {
	var buf bytes.Buffer
	printReportTo(&buf, nil, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Errorf("expected src/dst targets in header: %s", out)
	}
}
