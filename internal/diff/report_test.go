package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintReport_NoResults(t *testing.T) {
	var buf bytes.Buffer
	printReportTo(&buf, map[string]Result{})
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestPrintReport_OnlyInA(t *testing.T) {
	var buf bytes.Buffer
	result := map[string]Result{
		"FOO": {ValueA: "bar", OnlyInA: true},
	}
	printReportTo(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "only in A") {
		t.Errorf("expected 'only in A', got: %q", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output, got: %q", out)
	}
}

func TestPrintReport_OnlyInB(t *testing.T) {
	var buf bytes.Buffer
	result := map[string]Result{
		"BAR": {ValueB: "baz", OnlyInB: true},
	}
	printReportTo(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "only in B") {
		t.Errorf("expected 'only in B', got: %q", out)
	}
}

func TestPrintReport_Differing(t *testing.T) {
	var buf bytes.Buffer
	result := map[string]Result{
		"HOST": {ValueA: "localhost", ValueB: "prod.example.com"},
	}
	printReportTo(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "localhost") || !strings.Contains(out, "prod.example.com") {
		t.Errorf("expected both values in output, got: %q", out)
	}
}

func TestSummary_Identical(t *testing.T) {
	s := Summary(map[string]Result{})
	if s != "files are identical" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_Mixed(t *testing.T) {
	result := map[string]Result{
		"A": {ValueA: "x", OnlyInA: true},
		"B": {ValueB: "y", OnlyInB: true},
		"C": {ValueA: "1", ValueB: "2"},
	}
	s := Summary(result)
	if !strings.Contains(s, "1 only-in-A") {
		t.Errorf("expected 1 only-in-A in summary, got: %q", s)
	}
	if !strings.Contains(s, "1 only-in-B") {
		t.Errorf("expected 1 only-in-B in summary, got: %q", s)
	}
	if !strings.Contains(s, "1 differing") {
		t.Errorf("expected 1 differing in summary, got: %q", s)
	}
}
