package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// AppRoleChange describes a single field change between two AppRole snapshots.
type AppRoleChange struct {
	Field  string
	OldVal string
	NewVal string
}

// CompareAppRoles returns changes between two AppRoleInfo structs.
func CompareAppRoles(a, b *vault.AppRoleInfo) []AppRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []AppRoleChange

	if a.BindSecretID != b.BindSecretID {
		changes = append(changes, AppRoleChange{"bind_secret_id", fmt.Sprintf("%v", a.BindSecretID), fmt.Sprintf("%v", b.BindSecretID)})
	}
	if a.TTL != b.TTL {
		changes = append(changes, AppRoleChange{"token_ttl", fmt.Sprintf("%d", a.TTL), fmt.Sprintf("%d", b.TTL)})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, AppRoleChange{"token_max_ttl", fmt.Sprintf("%d", a.MaxTTL), fmt.Sprintf("%d", b.MaxTTL)})
	}

	aSet := sliceToStrSet(a.Policies)
	bSet := sliceToStrSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, AppRoleChange{"policy_added", "", p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, AppRoleChange{"policy_removed", p, ""})
		}
	}
	return changes
}

func sliceToStrSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// PrintAppRoleDiff prints the diff to stdout.
func PrintAppRoleDiff(changes []AppRoleChange) {
	FprintAppRoleDiff(os.Stdout, changes)
}

// FprintAppRoleDiff writes the diff to the given writer.
func FprintAppRoleDiff(w io.Writer, changes []AppRoleChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No AppRole changes detected.")
		return
	}
	fmt.Fprintln(w, "AppRole Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch {
		case c.OldVal == "":
			fmt.Fprintf(w, "  + %-20s %s\n", c.Field, c.NewVal)
		case c.NewVal == "":
			fmt.Fprintf(w, "  - %-20s %s\n", c.Field, c.OldVal)
		default:
			fmt.Fprintf(w, "  ~ %-20s %s -> %s\n", c.Field, c.OldVal, c.NewVal)
		}
	}
}
