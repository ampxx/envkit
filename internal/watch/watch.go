package watch

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// FileEvent represents a detected change in a watched file.
type FileEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watcher polls one or more files for changes.
type Watcher struct {
	files    []string
	interval time.Duration
	hashes   map[string]string
	Events   chan FileEvent
	stop     chan struct{}
}

// New creates a Watcher for the given files and poll interval.
func New(files []string, interval time.Duration) *Watcher {
	return &Watcher{
		files:    files,
		interval: interval,
		hashes:   make(map[string]string),
		Events:   make(chan FileEvent, 16),
		stop:     make(chan struct{}),
	}
}

// Start begins polling in a goroutine. Call Stop to terminate.
func (w *Watcher) Start() {
	// Seed initial hashes so first tick doesn't fire spurious events.
	for _, f := range w.files {
		if h, err := hashFile(f); err == nil {
			w.hashes[f] = h
		}
	}
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop terminates the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	for _, f := range w.files {
		newHash, err := hashFile(f)
		if err != nil {
			continue
		}
		oldHash := w.hashes[f]
		if newHash != oldHash {
			w.Events <- FileEvent{
				Path:    f,
				OldHash: oldHash,
				NewHash: newHash,
				At:      time.Now(),
			}
			w.hashes[f] = newHash
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
