package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// FilterOptions controls which entries are returned by Filter.
type FilterOptions struct {
	Path      string    // exact path match; empty means all
	User      string    // exact user match; empty means all
	Since     time.Time // zero means no lower bound
	Until     time.Time // zero means no upper bound
}

// Filter reads newline-delimited JSON entries from r and returns those
// that satisfy all non-zero fields in opts.
func Filter(r io.Reader, opts FilterOptions) ([]Entry, error) {
	var results []Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		if opts.Path != "" && e.Path != opts.Path {
			continue
		}
		if opts.User != "" && e.User != opts.User {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		results = append(results, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan: %w", err)
	}
	return results, nil
}
