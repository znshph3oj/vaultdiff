package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// UserpassChange represents a single field change in a userpass role.
type UserpassChange struct {
	Field    string
	OldValue string
	NewValue string
}

// CompareUserpassRoles returns a list of changes between two UserpassRoleInfo structs.
func CompareUserpassRoles(a, b *vault.UserpassRoleInfo) []UserpassChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []UserpassChange

	if a.TTL != b.TTL {
		changes = append(changes, UserpassChange{Field: "token_ttl", OldValue: a.TTL, NewValue: b.TTL})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, UserpassChange{Field: "token_max_ttl", OldValue: a.MaxTTL, NewValue: b.MaxTTL})
	}

	aSet := sliceToStrSet(a.Policies)
	bSet := sliceToStrSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, UserpassChange{Field: "policy_added", OldValue: "", NewValue: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, UserpassChange{Field: "policy_removed", OldValue: p, NewValue: ""})
		}
	}

	aCIDRs := strings.Join(a.BoundCIDRs, ",")
	bCIDRs := strings.Join(b.BoundCIDRs, ",")
	if aCIDRs != bCIDRs {
		changes = append(changes, UserpassChange{Field: "token_bound_cidrs", OldValue: aCIDRs, NewValue: bCIDRs})
	}

	return changes
}

// PrintUserpassDiff prints the diff to stdout.
func PrintUserpassDiff(username string, changes []UserpassChange) {
	FprintUserpassDiff(os.Stdout, username, changes)
}

// FprintUserpassDiff writes the diff to the given writer.
func FprintUserpassDiff(w io.Writer, username string, changes []UserpassChange) {
	if len(changes) == 0 {
		fmt.Fprintf(w, "userpass/%s: no changes\n", username)
		return
	}
	fmt.Fprintf(w, "userpass/%s:\n", username)
	for _, c := range changes {
		switch c.Field {
		case "policy_added":
			fmt.Fprintf(w, "  + policy: %s\n", c.NewValue)
		case "policy_removed":
			fmt.Fprintf(w, "  - policy: %s\n", c.OldValue)
		default:
			fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.OldValue, c.NewValue)
		}
	}
}
