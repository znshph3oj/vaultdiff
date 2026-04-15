package audit

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// Summary holds aggregated statistics over a slice of audit entries.
type Summary struct {
	TotalEvents  int
	ByUser       map[string]int
	ByPath       map[string]int
	ByOperation  map[string]int
}

// Summarize aggregates the provided entries into a Summary.
func Summarize(entries []Entry) Summary {
	s := Summary{
		ByUser:      make(map[string]int),
		ByPath:      make(map[string]int),
		ByOperation: make(map[string]int),
	}
	for _, e := range entries {
		s.TotalEvents++
		if e.User != "" {
			s.ByUser[e.User]++
		}
		if e.Path != "" {
			s.ByPath[e.Path]++
		}
		if e.Operation != "" {
			s.ByOperation[e.Operation]++
		}
	}
	return s
}

// PrintSummary writes a human-readable summary table to w.
// If w is nil, os.Stdout is used.
func PrintSummary(s Summary, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Total events:\t%d\n", s.TotalEvents)

	if len(s.ByUser) > 0 {
		fmt.Fprintln(tw, "\nEvents by user:")
		for _, u := range sortedKeys(s.ByUser) {
			fmt.Fprintf(tw, "  %s\t%d\n", u, s.ByUser[u])
		}
	}

	if len(s.ByPath) > 0 {
		fmt.Fprintln(tw, "\nEvents by path:")
		for _, p := range sortedKeys(s.ByPath) {
			fmt.Fprintf(tw, "  %s\t%d\n", p, s.ByPath[p])
		}
	}

	if len(s.ByOperation) > 0 {
		fmt.Fprintln(tw, "\nEvents by operation:")
		for _, op := range sortedKeys(s.ByOperation) {
			fmt.Fprintf(tw, "  %s\t%d\n", op, s.ByOperation[op])
		}
	}
	tw.Flush()
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
