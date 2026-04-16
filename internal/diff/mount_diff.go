package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// MountChange describes a change to a single mount path.
type MountChange struct {
	Path   string
	Status string // "added", "removed", "modified"
	Old    *vault.MountInfo
	New    *vault.MountInfo
}

// CompareMounts diffs two snapshots of Vault mount configurations.
func CompareMounts(before, after map[string]*vault.MountInfo) []MountChange {
	var changes []MountChange

	for path, newInfo := range after {
		if oldInfo, ok := before[path]; !ok {
			changes = append(changes, MountChange{Path: path, Status: "added", New: newInfo})
		} else if oldInfo.Type != newInfo.Type || oldInfo.Description != newInfo.Description ||
			oldInfo.Local != newInfo.Local || oldInfo.SealWrap != newInfo.SealWrap {
			changes = append(changes, MountChange{Path: path, Status: "modified", Old: oldInfo, New: newInfo})
		}
	}

	for path, oldInfo := range before {
		if _, ok := after[path]; !ok {
			changes = append(changes, MountChange{Path: path, Status: "removed", Old: oldInfo})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})
	return changes
}

// PrintMountDiff writes mount diff output to stdout.
func PrintMountDiff(changes []MountChange) {
	FprintMountDiff(os.Stdout, changes)
}

// FprintMountDiff writes mount diff output to the provided writer.
func FprintMountDiff(w io.Writer, changes []MountChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No mount changes detected.")
		return
	}
	fmt.Fprintln(w, "Mount Changes:")
	fmt.Fprintln(w, "──────────────")
	for _, c := range changes {
		switch c.Status {
		case "added":
			fmt.Fprintf(w, "  + %s (type: %s)\n", c.Path, c.New.Type)
		case "removed":
			fmt.Fprintf(w, "  - %s (type: %s)\n", c.Path, c.Old.Type)
		case "modified":
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Path, c.Old.Type, c.New.Type)
		}
	}
}
