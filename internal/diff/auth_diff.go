package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// AuthChange describes a single difference between two AuthInfo snapshots.
type AuthChange struct {
	Field  string
	Before string
	After  string
}

// CompareAuth returns a list of changes between two AuthInfo values.
func CompareAuth(a, b *vault.AuthInfo) []AuthChange {
	var changes []AuthChange

	if a == nil || b == nil {
		return changes
	}
	if a.DisplayName != b.DisplayName {
		changes = append(changes, AuthChange{Field: "display_name", Before: a.DisplayName, After: b.DisplayName})
	}
	if a.Renewable != b.Renewable {
		changes = append(changes, AuthChange{
			Field:  "renewable",
			Before: fmt.Sprintf("%v", a.Renewable),
			After:  fmt.Sprintf("%v", b.Renewable),
		})
	}
	if a.TTL != b.TTL {
		changes = append(changes, AuthChange{
			Field:  "ttl",
			Before: a.TTL.String(),
			After:  b.TTL.String(),
		})
	}

	aSet := make(map[string]bool)
	bSet := make(map[string]bool)
	for _, p := range a.Policies {
		aSet[p] = true
	}
	for _, p := range b.Policies {
		bSet[p] = true
	}
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, AuthChange{Field: "policy", Before: "", After: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, AuthChange{Field: "policy", Before: p, After: ""})
		}
	}

	return changes
}

// PrintAuthDiff prints auth differences to stdout.
func PrintAuthDiff(changes []AuthChange) {
	FprintAuthDiff(os.Stdout, changes)
}

// FprintAuthDiff writes auth differences to w.
func FprintAuthDiff(w io.Writer, changes []AuthChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No auth changes detected.")
		return
	}
	fmt.Fprintln(w, "Auth Changes:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch {
		case c.Before == "":
			fmt.Fprintf(w, "  + %-20s %s\n", c.Field, c.After)
		case c.After == "":
			fmt.Fprintf(w, "  - %-20s %s\n", c.Field, c.Before)
		default:
			fmt.Fprintf(w, "  ~ %-20s %s -> %s\n", c.Field, c.Before, c.After)
		}
	}
}
