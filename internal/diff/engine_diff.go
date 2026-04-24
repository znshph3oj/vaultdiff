package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// EngineDiffEntry represents a change between two sets of mounted engines.
type EngineDiffEntry struct {
	Path   string
	Status string // "added", "removed", "modified"
	Old    *vault.EngineInfo
	New    *vault.EngineInfo
}

// CompareEngines diffs two slices of EngineInfo by path.
func CompareEngines(before, after []vault.EngineInfo) []EngineDiffEntry {
	beforeMap := make(map[string]vault.EngineInfo, len(before))
	for _, e := range before {
		beforeMap[e.Path] = e
	}
	afterMap := make(map[string]vault.EngineInfo, len(after))
	for _, e := range after {
		afterMap[e.Path] = e
	}

	keys := make(map[string]struct{})
	for k := range beforeMap {
		keys[k] = struct{}{}
	}
	for k := range afterMap {
		keys[k] = struct{}{}
	}

	var entries []EngineDiffEntry
	for path := range keys {
		old, hadOld := beforeMap[path]
		new, hasNew := afterMap[path]
		switch {
		case hadOld && !hasNew:
			entries = append(entries, EngineDiffEntry{Path: path, Status: "removed", Old: &old})
		case !hadOld && hasNew:
			entries = append(entries, EngineDiffEntry{Path: path, Status: "added", New: &new})
		case old.Type != new.Type || old.Description != new.Description || old.Local != new.Local || old.SealWrap != new.SealWrap:
			entries = append(entries, EngineDiffEntry{Path: path, Status: "modified", Old: &old, New: &new})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	return entries
}

// FilterByStatus returns only the diff entries matching the given status.
// Valid statuses are "added", "removed", and "modified".
func FilterByStatus(entries []EngineDiffEntry, status string) []EngineDiffEntry {
	var filtered []EngineDiffEntry
	for _, e := range entries {
		if e.Status == status {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// Summary returns a short string summarising the counts of added, removed,
// and modified entries in the diff (e.g. "+2 -1 ~3").
func Summary(entries []EngineDiffEntry) string {
	var added, removed, modified int
	for _, e := range entries {
		switch e.Status {
		case "added":
			added++
		case "removed":
			removed++
		case "modified":
			modified++
		}
	}
	return fmt.Sprintf("+%d -%d ~%d", added, removed, modified)
}

// PrintEngineDiff writes a human-readable engine diff to stdout.
func PrintEngineDiff(entries []EngineDiffEntry) {
	FprintEngineDiff(os.Stdout, entries)
}

// FprintEngineDiff writes a human-readable engine diff to w.
func FprintEngineDiff(w io.Writer, entries []EngineDiffEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No engine changes detected.")
		return
	}
	fmt.Fprintln(w, "Engine Mount Diff:")
	for _, e := range entries {
		switch e.Status {
		case "added":
			fmt.Fprintf(w, "  + %s (type: %s)\n", e.Path, e.New.Type)
		case "removed":
			fmt.Fprintf(w, "  - %s (type: %s)\n", e.Path, e.Old.Type)
		case "modified":
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", e.Path, e.Old.Type, e.New.Type)
		}
	}
	fmt.Fprintf(w, "Summary: %s\n", Summary(entries))
}
