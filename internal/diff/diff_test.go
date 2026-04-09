package diff

import (
	"strings"
	"testing"
)

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "SHARED": "val"}
	b := map[string]string{"SHARED": "val"}

	r := Compare(a, b)

	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "FOO" {
		t.Errorf("expected OnlyInA=[FOO], got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 0 {
		t.Errorf("expected OnlyInB empty, got %v", r.OnlyInB)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"SHARED": "val"}
	b := map[string]string{"BAZ": "qux", "SHARED": "val"}

	r := Compare(a, b)

	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BAZ" {
		t.Errorf("expected OnlyInB=[BAZ], got %v", r.OnlyInB)
	}
}

func TestCompare_Differing(t *testing.T) {
	a := map[string]string{"KEY": "alpha"}
	b := map[string]string{"KEY": "beta"}

	r := Compare(a, b)

	if len(r.Differing) != 1 {
		t.Fatalf("expected 1 differing entry, got %d", len(r.Differing))
	}
	if r.Differing[0].Key != "KEY" || r.Differing[0].ValueA != "alpha" || r.Differing[0].ValueB != "beta" {
		t.Errorf("unexpected differing entry: %+v", r.Differing[0])
	}
}

func TestCompare_Common(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	r := Compare(a, b)

	if len(r.Common) != 1 || r.Common[0] != "KEY" {
		t.Errorf("expected Common=[KEY], got %v", r.Common)
	}
	if len(r.Differing) != 0 || len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Errorf("expected no differences")
	}
}

func TestFormat_NoDiff(t *testing.T) {
	r := Result{Common: []string{"KEY"}}
	out := Format(r, "staging", "production")
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-differences message, got: %s", out)
	}
}

func TestFormat_WithDiff(t *testing.T) {
	r := Result{
		OnlyInA:   []string{"ONLY_A"},
		OnlyInB:   []string{"ONLY_B"},
		Differing: []DiffEntry{{Key: "KEY", ValueA: "v1", ValueB: "v2"}},
	}
	out := Format(r, "staging", "production")
	if !strings.Contains(out, "ONLY_A") {
		t.Errorf("expected ONLY_A in output")
	}
	if !strings.Contains(out, "ONLY_B") {
		t.Errorf("expected ONLY_B in output")
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in output")
	}
}
