package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplates_ReturnsKnownNames(t *testing.T) {
	names := ListTemplates()
	if len(names) == 0 {
		t.Fatal("expected at least one template")
	}
	found := false
	for _, n := range names {
		if n == "web" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'web' template in list")
	}
}

func TestScaffold_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	err := Scaffold(Options{Dir: dir, Template: "web", Force: false})
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}
	dest := filepath.Join(dir, "envkit.yaml")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", dest)
	}
}

func TestScaffold_FileContainsTargets(t *testing.T) {
	dir := t.TempDir()
	_ = Scaffold(Options{Dir: dir, Template: "web"})
	data, err := os.ReadFile(filepath.Join(dir, "envkit.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, target := range []string{"development", "staging", "production"} {
		if !strings.Contains(content, target) {
			t.Errorf("expected target %q in output", target)
		}
	}
}

func TestScaffold_UnknownTemplate(t *testing.T) {
	dir := t.TempDir()
	err := Scaffold(Options{Dir: dir, Template: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
}

func TestScaffold_ErrorsIfFileExists(t *testing.T) {
	dir := t.TempDir()
	_ = Scaffold(Options{Dir: dir, Template: "web"})
	err := Scaffold(Options{Dir: dir, Template: "web", Force: false})
	if err == nil {
		t.Fatal("expected error when file already exists")
	}
}

func TestScaffold_ForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	_ = Scaffold(Options{Dir: dir, Template: "web"})
	err := Scaffold(Options{Dir: dir, Template: "worker", Force: true})
	if err != nil {
		t.Fatalf("expected no error with --force, got: %v", err)
	}
}

func TestScaffold_WorkerTemplate(t *testing.T) {
	dir := t.TempDir()
	err := Scaffold(Options{Dir: dir, Template: "worker"})
	if err != nil {
		t.Fatalf("Scaffold worker error: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "envkit.yaml"))
	if !strings.Contains(string(data), "QUEUE_URL") {
		t.Error("expected QUEUE_URL in worker template output")
	}
}
