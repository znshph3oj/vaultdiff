package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Path      string           `json:"path"`
	FromVersion int            `json:"from_version"`
	ToVersion   int            `json:"to_version"`
	Changes   []diff.Change    `json:"changes"`
	User      string           `json:"user,omitempty"`
}

// Logger writes audit entries to an output sink.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to w.
// Pass nil to use os.Stdout.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// NewFileLogger opens (or creates) a file at path and returns a Logger
// that appends audit entries to it.
func NewFileLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{w: f}, nil
}

// Record serialises an Entry as a single JSON line.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
