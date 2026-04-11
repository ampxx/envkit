package scope

import (
	"strings"
	"testing"
)

func TestPrintReport_ContainsTarget(t *testing.T) {
	r := Result{
		Target: "staging",
		Vars:   map[string]string{"HOST": "stg.example.com"},
	}
	var buf strings.Builder
	printReportTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected target name in output, got: %s", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected key HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "stg.example.com") {
		t.Errorf("expected value in output, got: %s", out)
	}
}

func TestPrintReport_ShowsVarCount(t *testing.T) {
	r := Result{
		Target: "prod",
		Vars:   map[string]string{"A": "1", "B": "2"},
	}
	var buf strings.Builder
	printReportTo(&buf, r)
	if !strings.Contains(buf.String(), "vars=2") {
		t.Errorf("expected vars=2 in output")
	}
}

func TestSummary_Format(t *testing.T) {
	r := Result{Target: "dev", Vars: map[string]string{"X": "y", "Z": "w"}}
	s := Summary(r)
	if !strings.Contains(s, "dev") {
		t.Errorf("expected target in summary")
	}
	if !strings.Contains(s, "2") {
		t.Errorf("expected count in summary")
	}
}
