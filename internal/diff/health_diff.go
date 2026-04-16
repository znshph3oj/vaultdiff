package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// HealthChange describes a change between two health snapshots.
type HealthChange struct {
	Field  string
	Before string
	After  string
}

// CompareHealth returns changes between two HealthStatus values.
func CompareHealth(a, b *vault.HealthStatus) []HealthChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []HealthChange
	check := func(field, before, after string) {
		if before != after {
			changes = append(changes, HealthChange{Field: field, Before: before, After: after})
		}
	}
	check("initialized", fmt.Sprintf("%v", a.Initialized), fmt.Sprintf("%v", b.Initialized))
	check("sealed", fmt.Sprintf("%v", a.Sealed), fmt.Sprintf("%v", b.Sealed))
	check("standby", fmt.Sprintf("%v", a.Standby), fmt.Sprintf("%v", b.Standby))
	check("version", a.Version, b.Version)
	check("cluster_name", a.ClusterName, b.ClusterName)
	check("cluster_id", a.ClusterID, b.ClusterID)
	return changes
}

// PrintHealthDiff writes health changes to stdout.
func PrintHealthDiff(changes []HealthChange) {
	FprintHealthDiff(os.Stdout, changes)
}

// FprintHealthDiff writes health changes to the given writer.
func FprintHealthDiff(w io.Writer, changes []HealthChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No health changes detected.")
		return
	}
	fmt.Fprintln(w, "Health Changes:")
	for _, c := range changes {
		fmt.Fprintf(w, "  %-20s %s -> %s\n", c.Field, c.Before, c.After)
	}
}
