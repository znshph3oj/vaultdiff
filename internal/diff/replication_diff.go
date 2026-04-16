package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// ReplicationChange describes a single changed field in replication status.
type ReplicationChange struct {
	Field string
	From  string
	To    string
}

// CompareReplication returns changes between two ReplicationStatus snapshots.
func CompareReplication(a, b *vault.ReplicationStatus) []ReplicationChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []ReplicationChange
	if a.DRMode != b.DRMode {
		changes = append(changes, ReplicationChange{Field: "dr_mode", From: a.DRMode, To: b.DRMode})
	}
	if a.PerformanceMode != b.PerformanceMode {
		changes = append(changes, ReplicationChange{Field: "performance_mode", From: a.PerformanceMode, To: b.PerformanceMode})
	}
	if a.Primary != b.Primary {
		changes = append(changes, ReplicationChange{
			Field: "primary",
			From:  fmt.Sprintf("%v", a.Primary),
			To:    fmt.Sprintf("%v", b.Primary),
		})
	}
	aSet := toSet(a.KnownSecondaries)
	bSet := toSet(b.KnownSecondaries)
	for k := range bSet {
		if !aSet[k] {
			changes = append(changes, ReplicationChange{Field: "known_secondaries", From: "", To: k})
		}
	}
	for k := range aSet {
		if !bSet[k] {
			changes = append(changes, ReplicationChange{Field: "known_secondaries", From: k, To: ""})
		}
	}
	return changes
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

// PrintReplicationDiff prints replication changes to stdout.
func PrintReplicationDiff(changes []ReplicationChange) {
	FprintReplicationDiff(os.Stdout, changes)
}

// FprintReplicationDiff writes replication changes to the given writer.
func FprintReplicationDiff(w io.Writer, changes []ReplicationChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No replication changes detected.")
		return
	}
	fmt.Fprintln(w, "Replication Status Changes:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		if c.From == "" {
			fmt.Fprintf(w, "  + %-24s %s\n", c.Field, c.To)
		} else if c.To == "" {
			fmt.Fprintf(w, "  - %-24s %s\n", c.Field, c.From)
		} else {
			fmt.Fprintf(w, "  ~ %-24s %s -> %s\n", c.Field, c.From, c.To)
		}
	}
}
