package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// NomadChange describes a single field-level change between two NomadRoleInfo values.
type NomadChange struct {
	Field string
	From  string
	To    string
}

// CompareNomadRoles returns the list of changes between two Nomad role configs.
func CompareNomadRoles(a, b *vault.NomadRoleInfo) []NomadChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []NomadChange

	if a.Type != b.Type {
		changes = append(changes, NomadChange{Field: "type", From: a.Type, To: b.Type})
	}
	if a.Lease != b.Lease {
		changes = append(changes, NomadChange{Field: "lease", From: a.Lease, To: b.Lease})
	}
	if fmt.Sprintf("%v", a.Global) != fmt.Sprintf("%v", b.Global) {
		changes = append(changes, NomadChange{
			Field: "global",
			From:  fmt.Sprintf("%v", a.Global),
			To:    fmt.Sprintf("%v", b.Global),
		})
	}

	aSet := sliceToStrSet(a.Policies)
	bSet := sliceToStrSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, NomadChange{Field: "policy", From: "", To: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, NomadChange{Field: "policy", From: p, To: ""})
		}
	}
	return changes
}

// PrintNomadDiff prints a Nomad role diff to stdout.
func PrintNomadDiff(a, b *vault.NomadRoleInfo) {
	FprintNomadDiff(os.Stdout, a, b)
}

// FprintNomadDiff writes a Nomad role diff to the given writer.
func FprintNomadDiff(w io.Writer, a, b *vault.NomadRoleInfo) {
	changes := CompareNomadRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes in Nomad role configuration.")
		return
	}
	fmt.Fprintln(w, "Nomad Role Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch {
		case c.From == "":
			fmt.Fprintf(w, "  + %s: %s\n", c.Field, c.To)
		case c.To == "":
			fmt.Fprintf(w, "  - %s: %s\n", c.Field, c.From)
		default:
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Field, c.From, c.To)
		}
	}
}
