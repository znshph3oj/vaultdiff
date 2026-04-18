package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// TransitKeyChange describes a change in a transit key field.
type TransitKeyChange struct {
	Field string
	Old   interface{}
	New   interface{}
}

// CompareTransitKeys returns a list of changes between two TransitKeyInfo values.
func CompareTransitKeys(a, b *vault.TransitKeyInfo) []TransitKeyChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []TransitKeyChange

	if a.Type != b.Type {
		changes = append(changes, TransitKeyChange{Field: "type", Old: a.Type, New: b.Type})
	}
	if a.DeletionAllowed != b.DeletionAllowed {
		changes = append(changes, TransitKeyChange{Field: "deletion_allowed", Old: a.DeletionAllowed, New: b.DeletionAllowed})
	}
	if a.Exportable != b.Exportable {
		changes = append(changes, TransitKeyChange{Field: "exportable", Old: a.Exportable, New: b.Exportable})
	}
	if a.LatestVersion != b.LatestVersion {
		changes = append(changes, TransitKeyChange{Field: "latest_version", Old: a.LatestVersion, New: b.LatestVersion})
	}
	if a.MinDecryptVersion != b.MinDecryptVersion {
		changes = append(changes, TransitKeyChange{Field: "min_decryption_version", Old: a.MinDecryptVersion, New: b.MinDecryptVersion})
	}
	return changes
}

// PrintTransitDiff prints transit key changes to stdout.
func PrintTransitDiff(keyName string, changes []TransitKeyChange) {
	FprintTransitDiff(os.Stdout, keyName, changes)
}

// FprintTransitDiff writes transit key diff to the given writer.
func FprintTransitDiff(w io.Writer, keyName string, changes []TransitKeyChange) {
	if len(changes) == 0 {
		fmt.Fprintf(w, "transit key %q: no changes\n", keyName)
		return
	}
	fmt.Fprintf(w, "transit key %q changes:\n", keyName)
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %v -> %v\n", c.Field, c.Old, c.New)
	}
}
