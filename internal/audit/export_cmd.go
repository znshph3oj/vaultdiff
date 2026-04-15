package audit

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ExportOptions configures an export operation.
type ExportOptions struct {
	Format   ExportFormat
	Output   string // file path; empty means stdout
	Since    time.Time
	User     string
	Path     string
}

// RunExport filters and exports audit entries according to opts.
// Entries are sourced from the provided slice (e.g. loaded from a log file).
func RunExport(entries []Entry, opts ExportOptions) error {
	filtered := filterEntries(entries, opts)

	var w io.Writer = os.Stdout
	if opts.Output != "" {
		f, err := os.Create(opts.Output)
		if err != nil {
			return fmt.Errorf("creating output file %q: %w", opts.Output, err)
		}
		defer f.Close()
		w = f
	}

	if err := Export(w, filtered, opts.Format); err != nil {
		return fmt.Errorf("exporting audit log: %w", err)
	}

	return nil
}

// filterEntries applies ExportOptions filters to a slice of entries.
func filterEntries(entries []Entry, opts ExportOptions) []Entry {
	var out []Entry
	for _, e := range entries {
		if opts.User != "" && e.User != opts.User {
			continue
		}
		if opts.Path != "" && e.Path != opts.Path {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		out = append(out, e)
	}
	return out
}
