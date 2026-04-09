package prompt

import (
	"bufio"
	"strings"
	"testing"
)

func newTestPrompter(input string) *Prompter {
	return &Prompter{reader: bufio.NewReader(strings.NewReader(input))}
}

func TestAskString_WithInput(t *testing.T) {
	p := newTestPrompter("myvalue\n")
	val, err := p.AskString("MY_KEY", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "myvalue" {
		t.Errorf("expected 'myvalue', got %q", val)
	}
}

func TestAskString_UsesDefault(t *testing.T) {
	p := newTestPrompter("\n")
	val, err := p.AskString("MY_KEY", "default_val")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "default_val" {
		t.Errorf("expected 'default_val', got %q", val)
	}
}

func TestAskConfirm_Yes(t *testing.T) {
	for _, input := range []string{"y\n", "yes\n", "Y\n", "YES\n"} {
		p := newTestPrompter(input)
		ok, err := p.AskConfirm("Continue?")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Errorf("expected true for input %q", input)
		}
	}
}

func TestAskConfirm_No(t *testing.T) {
	for _, input := range []string{"n\n", "no\n", "\n"} {
		p := newTestPrompter(input)
		ok, err := p.AskConfirm("Continue?")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Errorf("expected false for input %q", input)
		}
	}
}

func TestFillMissing(t *testing.T) {
	// Two keys, first filled, second skipped
	p := newTestPrompter("value1\n\n")
	result, err := p.FillMissing([]string{"KEY_ONE", "KEY_TWO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY_ONE"] != "value1" {
		t.Errorf("expected KEY_ONE=value1, got %q", result["KEY_ONE"])
	}
	if _, ok := result["KEY_TWO"]; ok {
		t.Error("expected KEY_TWO to be absent (skipped)")
	}
}

func TestFillMissing_Empty(t *testing.T) {
	p := newTestPrompter("")
	result, err := p.FillMissing(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for empty keys, got %v", result)
	}
}
