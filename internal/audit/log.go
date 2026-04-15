package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	User      string            `json:"user"`
	Operation string            `json:"operation"`
	Path      string            `json:"path"`
	Data      map[string]string `json:"data,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	w.Writer
}

// NewLogger returns a Logger writing to the given writer.
func NewLogger(w io.Writer) *Logger {
	return &Logger{w: NewFileLogger opens (or creates) a file at path and returns a Logger for it.
func NewFileLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}
	return NewLogger(f), nil
}

// Record writes the entry to the underlying writer as a JSON line.
// If Timestamp is zero it is set to the current UTC time.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	enc := json.NewEncoder(l.w)
	return enc.Encode(e)
}

// ReadEntries decodes newline-delimited JSON entries from r.
func ReadEntries(r io.Reader) ([]Entry, error) {
	var entries []Entry
	dec := json.NewDecoder(r)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
