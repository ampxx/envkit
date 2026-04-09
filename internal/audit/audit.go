package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventValidate EventType = "validate"
	EventExport   EventType = "export"
	EventFill     EventType = "fill"
	EventSnapshot EventType = "snapshot"
	EventLint     EventType = "lint"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Target    string    `json:"target"`
	User      string    `json:"user,omitempty"`
	Details   string    `json:"details,omitempty"`
	Success   bool      `json:"success"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// New creates a new Logger that appends to the given file path.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Log appends an Entry to the audit log file.
func (l *Logger) Log(event EventType, target string, success bool, details string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Target:    target,
		Success:   success,
		Details:   details,
		User:      currentUser(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	_, err = fmt.Fprintln(f, string(data))
	return err
}

// ReadAll reads and returns all entries from the audit log file.
func ReadAll(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read log: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(string(data)) {
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func currentUser() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return os.Getenv("USERNAME")
}
