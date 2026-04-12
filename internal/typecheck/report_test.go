package typecheck

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintReport_ShowsOK(t *testing.T) {
	r := Report{
		Target: "production",
		Results: []Result{
			{Key: "PORT", Expected: "int", Actual: "8080", Status: StatusOK},
		},
	}
	var buf bytes.Buffer
	printReportTo(&buf, r)
	if !strings.Contains(buf.String(), "✓") {
		t.Error("expected ✓ for OK result")
	}
	if !strings.Contains(buf.String(), "PORT") {
		t.Error("expected key name in output")
	}
}

func TestPrintReport_ShowsFail(t *testing.T) {
	r := Report{
		Target: "production",
		Results: []Result{
			{Key: "PORT", Expected: "int", Actual: "abc", Status: StatusFail, Reason: `"abc" is not an integer`},
		},
	}
	var buf bytes.Buffer
	printReportTo(&buf, r)
	if !strings.Contains(buf.String(), "✗") {
		t.Error("expected ✗ for failed result")
	}
	if !strings.Contains(buf.String(), "not an integer") {
		t.Error("expected reason in output")
	}
}

func TestPrintReport_ShowsSkipped(t *testing.T) {
	r := Report{
		Target: "production",
		Results: []Result{
			{Key: "FOO", Expected: "any", Status: StatusSkipped},
		},
	}
	var buf bytes.Buffer
	printReportTo(&buf, r)
	if !strings.Contains(buf.String(), "-") {
		t.Error("expected dash for skipped result")
	}
}

func TestSummary_Counts(t *testing.T) {
	r := Report{
		Results: []Result{
			{Status: StatusOK},
			{Status: StatusOK},
			{Status: StatusFail},
			{Status: StatusSkipped},
		},
	}
	s := Summary(r)
	if !strings.Contains(s, "2 ok") {
		t.Errorf("expected '2 ok' in %q", s)
	}
	if !strings.Contains(s, "1 failed") {
		t.Errorf("expected '1 failed' in %q", s)
	}
	if !strings.Contains(s, "1 skipped") {
		t.Errorf("expected '1 skipped' in %q", s)
	}
}
