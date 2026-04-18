package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// GCPRoleChange represents a field-level change in a GCP role.
type GCPRoleChange struct {
	Field string
	From  string
	To    string
}

// CompareGCPRoles compares two GCPRoleInfo structs and returns a list of changes.
func CompareGCPRoles(a, b *vault.GCPRoleInfo) []GCPRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []GCPRoleChange

	check := func(field, from, to string) {
		if from != to {
			changes = append(changes, GCPRoleChange{Field: field, From: from, To: to})
		}
	}

	check("role_type", a.RoleType, b.RoleType)
	check("project", a.Project, b.Project)
	check("secret_type", a.SecretType, b.SecretType)
	check("bindings", a.Bindings, b.Bindings)
	check("token_ttl", fmt.Sprintf("%d", a.TokenTTL), fmt.Sprintf("%d", b.TokenTTL))
	check("token_scopes", strings.Join(a.Scopes, ","), strings.Join(b.Scopes, ","))

	return changes
}

// PrintGCPDiff prints the GCP role diff to stdout.
func PrintGCPDiff(a, b *vault.GCPRoleInfo) {
	FprintGCPDiff(os.Stdout, a, b)
}

// FprintGCPDiff writes the GCP role diff to the given writer.
func FprintGCPDiff(w io.Writer, a, b *vault.GCPRoleInfo) {
	changes := CompareGCPRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes in GCP role.")
		return
	}
	fmt.Fprintln(w, "GCP Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
