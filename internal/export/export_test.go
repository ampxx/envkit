package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envkit/envkit/internal/config"
)

func makeConfig() *config.Config {
	return &config.Config{
		Targets: []config.Target{
			{
				Name: "production",
				Vars: []config.VarSpec{
					{Key: "APP_ENV"},
					{Key: "DB_URL"},
				},
			},
		},
	}
}

var sampleEnv = map[string]string{
	"APP_ENV": "production",
	"DB_URL":  "postgres://localhost/mydb",
}

func TestExport_Dotenv(t *testing.T) {
	cfg := makeConfig()
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	err := Export(cfg, sampleEnv, Options{Format: FormatDotenv, Target: "production"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line, got: %s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	cfg := makeConfig()
	tmp := filepath.Join(t.TempDir(), "out.sh")

	err := Export(cfg, sampleEnv, Options{Format: FormatShell, Target: "production", OutFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "export APP_ENV=") {
		t.Errorf("expected shell export, got: %s", data)
	}
}

func TestExport_JSON(t *testing.T) {
	cfg := makeConfig()
	tmp := filepath.Join(t.TempDir(), "out.json")

	err := Export(cfg, sampleEnv, Options{Format: FormatJSON, Target: "production", OutFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), `"APP_ENV"`) {
		t.Errorf("expected JSON key, got: %s", data)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	cfg := makeConfig()
	err := Export(cfg, sampleEnv, Options{Format: "xml", Target: "production"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExport_UnknownTarget(t *testing.T) {
	cfg := makeConfig()
	err := Export(cfg, sampleEnv, Options{Format: FormatDotenv, Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}
