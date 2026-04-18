package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultdiff/internal/vault"
)

// PKICertChange represents a change between two PKI cert states.
type PKICertChange struct {
	Field    string
	OldValue string
	NewValue string
}

// ComparePKICerts compares two PKICertInfo structs and returns a list of changes.
func ComparePKICerts(a, b *vault.PKICertInfo) []PKICertChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.PKICertInfo{}
	}
	if b == nil {
		b = &vault.PKICertInfo{}
	}

	var changes []PKICertChange
	check := func(field, ov, nv string) {
		if ov != nv {
			changes = append(changes, PKICertChange{Field: field, OldValue: ov, NewValue: nv})
		}
	}

	check("CommonName", a.CommonName, b.CommonName)
	check("Issuer", a.Issuer, b.Issuer)
	check("NotBefore", a.NotBefore, b.NotBefore)
	check("NotAfter", a.NotAfter, b.NotAfter)

	oldRevoked := fmt.Sprintf("%v", a.Revoked)
	newRevoked := fmt.Sprintf("%v", b.Revoked)
	check("Revoked", oldRevoked, newRevoked)

	return changes
}

// PrintPKIDiff prints PKI cert diff to stdout.
func PrintPKIDiff(serial string, changes []PKICertChange) {
	FprintPKIDiff(os.Stdout, serial, changes)
}

// FprintPKIDiff writes PKI cert diff to the given writer.
func FprintPKIDiff(w io.Writer, serial string, changes []PKICertChange) {
	if len(changes) == 0 {
		fmt.Fprintf(w, "PKI cert %s: no changes\n", serial)
		return
	}
	fmt.Fprintf(w, "PKI cert diff for %s:\n", serial)
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.OldValue, c.NewValue)
	}
}
