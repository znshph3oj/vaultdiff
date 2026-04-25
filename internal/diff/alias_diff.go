package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// AliasDiffEntry records a single field change between two AliasInfo values.
type AliasDiffEntry struct {
	Field string
	Old   string
	New   string
}

// CompareAliases returns the list of field-level differences between two AliasInfo structs.
func CompareAliases(a, b *vault.AliasInfo) []AliasDiffEntry {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.AliasInfo{}
	}
	if b == nil {
		b = &vault.AliasInfo{}
	}

	var changes []AliasDiffEntry

	check := func(field, oldVal, newVal string) {
		if oldVal != newVal {
			changes = append(changes, AliasDiffEntry{Field: field, Old: oldVal, New: newVal})
		}
	}

	check("name", a.Name, b.Name)
	check("mount_accessor", a.MountAccessor, b.MountAccessor)
	check("mount_type", a.MountType, b.MountType)
	check("canonical_id", a.CanonicalID, b.CanonicalID)

	// Compare metadata keys
	allKeys := unionStringKeys(a.Metadata, b.Metadata)
	for _, k := range allKeys {
		check(fmt.Sprintf("metadata.%s", k), a.Metadata[k], b.Metadata[k])
	}

	return changes
}

// PrintAliasDiff writes alias diff results to stdout.
func PrintAliasDiff(changes []AliasDiffEntry) {
	FprintAliasDiff(os.Stdout, changes)
}

// FprintAliasDiff writes alias diff results to the given writer.
func FprintAliasDiff(w io.Writer, changes []AliasDiffEntry) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No alias changes detected.")
		return
	}
	fmt.Fprintln(w, "Alias Changes:")
	for _, c := range changes {
		fmt.Fprintf(w, "  %-24s %s -> %s\n", c.Field+":", c.Old, c.New)
	}
}
