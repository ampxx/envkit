package rotate

import (
	"strings"
	"testing"
)

func TestPrintReport_ShowsRotated(t *testing.T) {
	r := Result{
		Rotated: []RotatedVar{
			{Key: "API_KEY", OldValue: "old-secret", NewValue: "new-secret"},
		},
	}
	var buf strings.Builder
	PrintReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "rotated:") {
		t.Errorf("expected 'rotated:' header, got: %s", out)
	}
}

func TestPrintReport_ShowsSkipped(t *testing.T) {
	r := Result{
		Skipped: []string{"DB_PASS"},
	}
	var buf strings.Builder
	PrintReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
}

func TestPrintReport_EmptyResult(t *testing.T) {
	r := Result{}
	var buf strings.Builder
	PrintReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "no variables matched") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestSummary_Counts(t *testing.T) {
	r := Result{
		Rotated: []RotatedVar{{Key: "A"}, {Key: "B"}},
		Skipped: []string{"C"},
	}
	s := Summary(r)
	if !strings.Contains(s, "2 rotated") {
		t.Errorf("expected '2 rotated' in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 skipped") {
		t.Errorf("expected '1 skipped' in summary, got: %s", s)
	}
}
