package watch

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/envkit/internal/audit"
	"github.com/user/envkit/internal/config"
	"github.com/user/envkit/internal/validator"
)

// RunWatch watches env files for a target and re-validates on change.
func RunWatch(cfgPath, targetName string, interval time.Duration) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var target *config.Target
	for i := range cfg.Targets {
		if cfg.Targets[i].Name == targetName {
			target = &cfg.Targets[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("target %q not found", targetName)
	}

	files := []string{target.EnvFile}
	w := New(files, interval)
	w.Start()
	defer w.Stop()

	fmt.Fprintf(os.Stdout, "Watching %s (target: %s) — press Ctrl+C to stop\n",
		target.EnvFile, targetName)

	// Run initial validation.
	runValidation(cfg, target, cfgPath)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev := <-w.Events:
			fmt.Fprintf(os.Stdout, "\n[changed] %s at %s\n", ev.Path, ev.At.Format("15:04:05"))
			runValidation(cfg, target, cfgPath)
		case <-sig:
			fmt.Fprintln(os.Stdout, "\nStopped watching.")
			return nil
		}
	}
}

func runValidation(cfg *config.Config, target *config.Target, cfgPath string) {
	results, err := validator.Validate(cfg, target.Name, target.EnvFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation error: %v\n", err)
		return
	}
	validator.PrintReport(os.Stdout, results)

	status := "pass"
	if validator.HasFailures(results) {
		status = "fail"
	}
	_ = audit.RecordEvent(cfgPath, "watch:validate",
		fmt.Sprintf("target=%s status=%s", target.Name, status))
}
