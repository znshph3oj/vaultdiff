package diff

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// AccessEntry records whether a user has read/write access to a path.
type AccessEntry struct {
	Path         string
	Capabilities []string
	CanRead      bool
	CanWrite     bool
}

// AccessReport holds access check results for one or more paths.
type AccessReport struct {
	User    string
	Entries []AccessEntry
}

// NewAccessReport builds an AccessReport from a map of path -> capabilities.
func NewAccessReport(user string, pathCaps map[string][]string, canRead func([]string) bool, canWrite func([]string) bool) *AccessReport {
	report := &AccessReport{User: user}
	for path, caps := range pathCaps {
		report.Entries = append(report.Entries, AccessEntry{
			Path:         path,
			Capabilities: caps,
			CanRead:      canRead(caps),
			CanWrite:     canWrite(caps),
		})
	}
	return report
}

// PrintAccessReport writes a formatted table of access entries to stdout.
func PrintAccessReport(r *AccessReport) {
	FprintAccessReport(os.Stdout, r)
}

// FprintAccessReport writes a formatted table of access entries to w.
func FprintAccessReport(w io.Writer, r *AccessReport) {
	fmt.Fprintf(w, "Access report for: %s\n", r.User)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tCAPABILITIES\tREAD\tWRITE")
	fmt.Fprintln(tw, strings.Repeat("-", 60))
	for _, e := range r.Entries {
		capsStr := strings.Join(e.Capabilities, ",")
		fmt.Fprintf(tw, "%s\t%s\t%v\t%v\n", e.Path, capsStr, e.CanRead, e.CanWrite)
	}
	tw.Flush()
}
