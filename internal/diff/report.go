package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Report holds the result of comparing two secret versions.
type Report struct {
	Path        string       `json:"path"`
	FromVersion int          `json:"from_version"`
	ToVersion   int          `json:"to_version"`
	GeneratedAt time.Time    `json:"generated_at"`
	Added       int          `json:"added"`
	Removed     int          `json:"removed"`
	Modified    int          `json:"modified"`
	Unchanged   int          `json:"unchanged"`
	Changes     []DiffResult `json:"changes"`
}

// NewReport builds a Report from a slice of DiffResults.
func NewReport(path string, from, to int, changes []DiffResult) *Report {
	r := &Report{
		Path:        path,
		FromVersion: from,
		ToVersion:   to,
		GeneratedAt: time.Now().UTC(),
		Changes:     changes,
	}
	for _, c := range changes {
		switch c.Status {
		case StatusAdded:
			r.Added++
		case StatusRemoved:
			r.Removed++
		case StatusModified:
			r.Modified++
		case StatusUnchanged:
			r.Unchanged++
		}
	}
	return r
}

// WriteJSON serialises the report as JSON to w.
func (r *Report) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// PrintSummary writes a human-readable summary to w.
func (r *Report) PrintSummary(w io.Writer) {
	fmt.Fprintf(w, "Path: %s  (v%d → v%d)\n", r.Path, r.FromVersion, r.ToVersion)
	fmt.Fprintf(w, "Added: %d  Removed: %d  Modified: %d  Unchanged: %d\n",
		r.Added, r.Removed, r.Modified, r.Unchanged)
}
