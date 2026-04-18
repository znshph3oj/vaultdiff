package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// IdentityChange represents a single field-level difference between two identity entities.
type IdentityChange struct {
	Field    string
	OldValue string
	NewValue string
}

// CompareIdentity returns a list of field-level changes between two IdentityEntity values.
// Returns nil if either input is nil.
func CompareIdentity(a, b *vault.IdentityEntity) []IdentityChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []IdentityChange

	if a.Name != b.Name {
		changes = append(changes, IdentityChange{Field: "name", OldValue: a.Name, NewValue: b.Name})
	}
	if a.Disabled != b.Disabled {
		changes = append(changes, IdentityChange{
			Field:    "disabled",
			OldValue: fmt.Sprintf("%v", a.Disabled),
			NewValue: fmt.Sprintf("%v", b.Disabled),
		})
	}

	aSet := sliceToSet(a.Policies)
	bSet := sliceToSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, IdentityChange{Field: "policy", OldValue: "", NewValue: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, IdentityChange{Field: "policy", OldValue: p, NewValue: ""})
		}
	}
	return changes
}

// sliceToSet converts a string slice into a set represented as a map.
func sliceToSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// PrintIdentityDiff prints a human-readable diff of two identity entities to stdout.
func PrintIdentityDiff(a, b *vault.IdentityEntity) {
	FprintIdentityDiff(os.Stdout, a, b)
}

// FprintIdentityDiff writes a human-readable diff of two identity entities to w.
func FprintIdentityDiff(w io.Writer, a, b *vault.IdentityEntity) {
	changes := CompareIdentity(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No identity changes detected.")
		return
	}
	fmt.Fprintln(w, "Identity Entity Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch {
		case c.OldValue == "":
			fmt.Fprintf(w, "  + [%s] %s\n", c.Field, c.NewValue)
		case c.NewValue == "":
			fmt.Fprintf(w, "  - [%s] %s\n", c.Field, c.OldValue)
		default:
			fmt.Fprintf(w, "  ~ [%s] %s -> %s\n", c.Field, c.OldValue, c.NewValue)
		}
	}
}
