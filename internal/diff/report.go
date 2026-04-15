package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Report represents a structured diff report between two secret versions.
type Report struct {
	Path      string       `json:"path"`
	VersionA  int          `json:"version_a"`
	VersionB  int          `json:"version_b"`
	Timestamp time.Time    `json:"timestamp"`
	Changes   []DiffResult `json:"changes"`
	Summary   ReportSummary `json:"summary"`
}

// ReportSummary holds counts of change types.
type ReportSummary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Unchanged int `json:"unchanged"`
}

// NewReport builds a Report from a slice of DiffResults.
func NewReport(path string, versionA, versionB int, results []DiffResult) Report {
	var summary ReportSummary
	for _, r := range results {
		switch r.Status {
		case StatusAdded:
			summary.Added++
		case StatusRemoved:
			summary.Removed++
		case StatusModified:
			summary.Modified++
		case StatusUnchanged:
			summary.Unchanged++
		}
	}
	return Report{
		Path:      path,
		VersionA:  versionA,
		VersionB:  versionB,
		Timestamp: time.Now().UTC(),
		Changes:   results,
		Summary:   summary,
	}
}

// WriteJSON serializes the report as JSON to the given writer.
func (r Report) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// PrintSummary writes a human-readable summary line to the given writer.
func (r Report) PrintSummary(w io.Writer) {
	fmt.Fprintf(w, "Diff report for %s (v%d → v%d): +%d added, -%d removed, ~%d modified, %d unchanged\n",
		r.Path, r.VersionA, r.VersionB,
		r.Summary.Added, r.Summary.Removed, r.Summary.Modified, r.Summary.Unchanged)
}
