package audit

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// PathStat holds aggregated change counts for a single secret path.
type PathStat struct {
	Path     string
	Added    int
	Removed  int
	Modified int
	Diffs    int // total diff operations recorded
}

// Summarize aggregates entries by path and returns stats sorted by path.
func Summarize(entries []Entry) []PathStat {
	m := make(map[string]*PathStat)
	for _, e := range entries {
		s, ok := m[e.Path]
		if !ok {
			s = &PathStat{Path: e.Path}
			m[e.Path] = s
		}
		s.Diffs++
		for _, c := range e.Changes {
			switch c.Type {
			case diff.Added:
				s.Added++
			case diff.Removed:
				s.Removed++
			case diff.Modified:
				s.Modified++
			}
		}
	}

	stats := make([]PathStat, 0, len(m))
	for _, s := range m {
		stats = append(stats, *s)
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Path < stats[j].Path
	})
	return stats
}

// PrintSummary writes a human-readable table of PathStats to w.
func PrintSummary(w io.Writer, stats []PathStat) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tDIFFS\tADDED\tREMOVED\tMODIFIED")
	for _, s := range stats {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\n",
			s.Path, s.Diffs, s.Added, s.Removed, s.Modified)
	}
	tw.Flush()
}
