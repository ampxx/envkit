package pin

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"envkit/internal/config"
)

// RunPin is the CLI entry point for the pin command.
func RunPin(cfgPath, target, pinnedBy string, keys []string, ttlStr string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	var ttl time.Duration
	if ttlStr != "" {
		ttl, err = time.ParseDuration(ttlStr)
		if err != nil {
			return fmt.Errorf("invalid ttl %q: %w", ttlStr, err)
		}
	}

	if pinnedBy == "" {
		pinnedBy = "unknown"
	}

	res, err := Pin(cfg, target, keys, pinnedBy, ttl)
	if err != nil {
		return err
	}

	printResult(w, res)
	return nil
}

// RunExpired prints entries that have passed their TTL.
func RunExpired(entries []Entry, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	exp := Expired(entries)
	if len(exp) == 0 {
		fmt.Fprintln(w, "no expired pins")
		return
	}
	fmt.Fprintf(w, "expired pins (%d):\n", len(exp))
	for _, e := range exp {
		fmt.Fprintf(w, "  [%s] %s (expired %s)\n", e.Target, e.Key, e.ExpiresAt.Format(time.RFC3339))
	}
}

func printResult(w io.Writer, res Result) {
	if len(res.Pinned) > 0 {
		fmt.Fprintf(w, "pinned (%d):\n", len(res.Pinned))
		for _, e := range res.Pinned {
			ttlNote := ""
			if !e.ExpiresAt.IsZero() {
				ttlNote = fmt.Sprintf(" [expires %s]", e.ExpiresAt.Format(time.RFC3339))
			}
			fmt.Fprintf(w, "  %s=%s%s\n", e.Key, e.Value, ttlNote)
		}
	}
	if len(res.Skipped) > 0 {
		fmt.Fprintf(w, "skipped: %s\n", strings.Join(res.Skipped, ", "))
	}
	if len(res.Errors) > 0 {
		fmt.Fprintf(w, "errors (%d):\n", len(res.Errors))
		for _, e := range res.Errors {
			fmt.Fprintf(w, "  ! %s\n", e)
		}
	}
	fmt.Fprintf(w, "summary: %d pinned, %d skipped, %d errors\n",
		len(res.Pinned), len(res.Skipped), len(res.Errors))
}
