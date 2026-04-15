package audit

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Path      string    `json:"path"`
	Operation string    `json:"operation"`
	VersionA  int       `json:"version_a,omitempty"`
	VersionB  int       `json:"version_b,omitempty"`
}

// ExportFormat defines the supported export formats.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// Export writes audit entries to w in the specified format.
func Export(w io.Writer, entries []Entry, format ExportFormat) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, entries)
	case FormatCSV:
		return exportCSV(w, entries)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func exportCSV(w io.Writer, entries []Entry) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"timestamp", "user", "path", "operation", "version_a", "version_b"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("writing csv header: %w", err)
	}

	for _, e := range entries {
		row := []string{
			e.Timestamp.UTC().Format(time.RFC3339),
			e.User,
			e.Path,
			e.Operation,
			fmt.Sprintf("%d", e.VersionA),
			fmt.Sprintf("%d", e.VersionB),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("writing csv row: %w", err)
		}
	}

	return cw.Error()
}
