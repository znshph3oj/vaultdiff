package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// RADIUSChange represents a single field change between two RADIUS role configs.
type RADIUSChange struct {
	Field string
	From  string
	To    string
}

// CompareRADIUSRoles compares two RADIUSRoleInfo structs and returns a list of changes.
func CompareRADIUSRoles(a, b *vault.RADIUSRoleInfo) []RADIUSChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.RADIUSRoleInfo{}
	}
	if b == nil {
		b = &vault.RADIUSRoleInfo{}
	}

	var changes []RADIUSChange

	if a.TTL != b.TTL {
		changes = append(changes, RADIUSChange{Field: "ttl", From: a.TTL, To: b.TTL})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, RADIUSChange{Field: "max_ttl", From: a.MaxTTL, To: b.MaxTTL})
	}

	aSet := sliceToStrSet(a.Policies)
	bSet := sliceToStrSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, RADIUSChange{Field: "policy_added", From: "", To: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, RADIUSChange{Field: "policy_removed", From: p, To: ""})
		}
	}

	return changes
}

// PrintRADIUSDiff prints the RADIUS role diff to stdout.
func PrintRADIUSDiff(changes []RADIUSChange) {
	FprintRADIUSDiff(os.Stdout, changes)
}

// FprintRADIUSDiff writes the RADIUS role diff to the given writer.
func FprintRADIUSDiff(w io.Writer, changes []RADIUSChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No RADIUS role changes detected.")
		return
	}
	fmt.Fprintln(w, "RADIUS Role Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch c.Field {
		case "policy_added":
			fmt.Fprintf(w, "  [+] policy: %s\n", c.To)
		case "policy_removed":
			fmt.Fprintf(w, "  [-] policy: %s\n", c.From)
		default:
			fmt.Fprintf(w, "  [~] %s: %q -> %q\n", c.Field, c.From, c.To)
		}
	}
}
