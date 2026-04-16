package diff

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/vaultdiff/internal/vault"
)

// LeaseChange represents a detected change in a lease's state.
type LeaseChange struct {
	LeaseID      string
	Field        string
	OldValue     interface{}
	NewValue     interface{}
}

// CompareLeases diffs two LeaseInfo snapshots and returns a list of changes.
// Returns nil if either snapshot is nil.
func CompareLeases(old, current *vault.LeaseInfo) []LeaseChange {
	if old == nil || current == nil {
		return nil
	}

	var changes []LeaseChange

	if old.Renewable != current.Renewable {
		changes = append(changes, LeaseChange{
			LeaseID:  current.LeaseID,
			Field:    "renewable",
			OldValue: old.Renewable,
			NewValue: current.Renewable,
		})
	}

	if old.LeaseDuration != current.LeaseDuration {
		changes = append(changes, LeaseChange{
			LeaseID:  current.LeaseID,
			Field:    "lease_duration",
			OldValue: old.LeaseDuration.String(),
			NewValue: current.LeaseDuration.String(),
		})
	}

	if !old.ExpireTime.Equal(current.ExpireTime) {
		changes = append(changes, LeaseChange{
			LeaseID:  current.LeaseID,
			Field:    "expire_time",
			OldValue: formatTime(old.ExpireTime),
			NewValue: formatTime(current.ExpireTime),
		})
	}

	return changes
}

// PrintLeaseDiff writes a human-readable lease diff to stdout.
func PrintLeaseDiff(changes []LeaseChange) {
	FprintLeaseDiff(os.Stdout, changes)
}

// FprintLeaseDiff writes a human-readable lease diff to the given writer.
func FprintLeaseDiff(w io.Writer, changes []LeaseChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No lease changes detected.")
		return
	}
	fmt.Fprintln(w, "Lease Changes:")
	for _, c := range changes {
		fmt.Fprintf(w, "  [%s] %s: %v -> %v\n", c.LeaseID, c.Field, c.OldValue, c.NewValue)
	}
}

// FilterChangesByField returns only the changes that match the given field name.
func FilterChangesByField(changes []LeaseChange, field string) []LeaseChange {
	var filtered []LeaseChange
	for _, c := range changes {
		if c.Field == field {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "(none)"
	}
	return t.UTC().Format(time.RFC3339)
}
