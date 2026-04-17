package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// RaftChange describes a change in Raft cluster membership or state.
type RaftChange struct {
	Field  string
	Before string
	After  string
}

// CompareRaft compares two RaftStatus snapshots and returns a list of changes.
func CompareRaft(a, b *vault.RaftStatus) []RaftChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []RaftChange

	if a.LeaderID != b.LeaderID {
		changes = append(changes, RaftChange{Field: "leader_id", Before: a.LeaderID, After: b.LeaderID})
	}
	if a.AppliedIndex != b.AppliedIndex {
		changes = append(changes, RaftChange{
			Field:  "applied_index",
			Before: fmt.Sprintf("%d", a.AppliedIndex),
			After:  fmt.Sprintf("%d", b.AppliedIndex),
		})
	}

	aNodes := make(map[string]vault.RaftServer)
	for _, s := range a.Servers {
		aNodes[s.NodeID] = s
	}
	bNodes := make(map[string]vault.RaftServer)
	for _, s := range b.Servers {
		bNodes[s.NodeID] = s
	}

	for id := range bNodes {
		if _, ok := aNodes[id]; !ok {
			changes = append(changes, RaftChange{Field: "server", Before: "", After: id + " (added)"})
		}
	}
	for id := range aNodes {
		if _, ok := bNodes[id]; !ok {
			changes = append(changes, RaftChange{Field: "server", Before: id + " (removed)", After: ""})
		}
	}

	return changes
}

// PrintRaftDiff prints Raft diff to stdout.
func PrintRaftDiff(changes []RaftChange) {
	FprintRaftDiff(os.Stdout, changes)
}

// FprintRaftDiff writes Raft diff to the given writer.
func FprintRaftDiff(w io.Writer, changes []RaftChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No raft changes detected.")
		return
	}
	fmt.Fprintln(w, "Raft Changes:")
	for _, c := range changes {
		if c.Before == "" {
			fmt.Fprintf(w, "  + [%s] %s\n", c.Field, c.After)
		} else if c.After == "" {
			fmt.Fprintf(w, "  - [%s] %s\n", c.Field, c.Before)
		} else {
			fmt.Fprintf(w, "  ~ [%s] %s -> %s\n", c.Field, c.Before, c.After)
		}
	}
}
