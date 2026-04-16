package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// AuditDeviceChange describes a change to an audit device.
type AuditDeviceChange struct {
	Path      string
	ChangeType string // added, removed, modified
	Old       *vault.AuditDevice
	New       *vault.AuditDevice
}

// CompareAuditDevices returns changes between two sets of audit devices.
func CompareAuditDevices(before, after map[string]*vault.AuditDevice) []AuditDeviceChange {
	var changes []AuditDeviceChange
	for path, dev := range after {
		if old, ok := before[path]; !ok {
			changes = append(changes, AuditDeviceChange{Path: path, ChangeType: "added", New: dev})
		} else if old.Type != dev.Type || old.Description != dev.Description {
			changes = append(changes, AuditDeviceChange{Path: path, ChangeType: "modified", Old: old, New: dev})
		}
	}
	for path, dev := range before {
		if _, ok := after[path]; !ok {
			changes = append(changes, AuditDeviceChange{Path: path, ChangeType: "removed", Old: dev})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return changes
}

// PrintAuditDeviceDiff prints audit device changes to stdout.
func PrintAuditDeviceDiff(changes []AuditDeviceChange) {
	FprintAuditDeviceDiff(os.Stdout, changes)
}

// FprintAuditDeviceDiff writes audit device changes to w.
func FprintAuditDeviceDiff(w io.Writer, changes []AuditDeviceChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No audit device changes.")
		return
	}
	fmt.Fprintln(w, "Audit Device Changes:")
	for _, c := range changes {
		switch c.ChangeType {
		case "added":
			fmt.Fprintf(w, "  + %s (type: %s)\n", c.Path, c.New.Type)
		case "removed":
			fmt.Fprintf(w, "  - %s (type: %s)\n", c.Path, c.Old.Type)
		case "modified":
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Path, c.Old.Type, c.New.Type)
		}
	}
}
