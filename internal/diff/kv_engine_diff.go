package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// KVEngineChange represents a detected change between two KV engine infos.
type KVEngineChange struct {
	Field    string
	OldValue string
	NewValue string
}

// CompareKVEngines compares two KVEngineInfo structs and returns changes.
func CompareKVEngines(a, b *vault.KVEngineInfo) []KVEngineChange {
	if a == nil && b == nil {
		return nil
	}
	var changes []KVEngineChange
	if a == nil || b == nil {
		changes = append(changes, KVEngineChange{Field: "engine", OldValue: fmt.Sprintf("%v", a), NewValue: fmt.Sprintf("%v", b)})
		return changes
	}
	if a.Version != b.Version {
		changes = append(changes, KVEngineChange{
			Field:    "version",
			OldValue: fmt.Sprintf("%d", a.Version),
			NewValue: fmt.Sprintf("%d", b.Version),
		})
	}
	allKeys := map[string]struct{}{}
	for k := range a.Options {
		allKeys[k] = struct{}{}
	}
	for k := range b.Options {
		allKeys[k] = struct{}{}
	}
	for k := range allKeys {
		av, bv := a.Options[k], b.Options[k]
		if av != bv {
			changes = append(changes, KVEngineChange{Field: "option:" + k, OldValue: av, NewValue: bv})
		}
	}
	return changes
}

// PrintKVEngineDiff prints the diff to stdout.
func PrintKVEngineDiff(changes []KVEngineChange) {
	FprintKVEngineDiff(os.Stdout, changes)
}

// FprintKVEngineDiff writes the diff to the given writer.
func FprintKVEngineDiff(w io.Writer, changes []KVEngineChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No KV engine changes detected.")
		return
	}
	fmt.Fprintln(w, "KV Engine Changes:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.OldValue, c.NewValue)
	}
}
