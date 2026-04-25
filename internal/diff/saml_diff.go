package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// SAMLRoleChange describes a single field change in a SAML role.
type SAMLRoleChange struct {
	Field string
	From  string
	To    string
}

// CompareSAMLRoles compares two SAMLRoleInfo structs and returns a list of changes.
func CompareSAMLRoles(a, b *vault.SAMLRoleInfo) []SAMLRoleChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.SAMLRoleInfo{}
	}
	if b == nil {
		b = &vault.SAMLRoleInfo{}
	}

	var changes []SAMLRoleChange

	if a.TokenTTL != b.TokenTTL {
		changes = append(changes, SAMLRoleChange{
			Field: "token_ttl",
			From:  fmt.Sprintf("%d", a.TokenTTL),
			To:    fmt.Sprintf("%d", b.TokenTTL),
		})
	}
	if a.TokenMaxTTL != b.TokenMaxTTL {
		changes = append(changes, SAMLRoleChange{
			Field: "token_max_ttl",
			From:  fmt.Sprintf("%d", a.TokenMaxTTL),
			To:    fmt.Sprintf("%d", b.TokenMaxTTL),
		})
	}

	aSubjects := strings.Join(a.BoundSubjects, ",")
	bSubjects := strings.Join(b.BoundSubjects, ",")
	if aSubjects != bSubjects {
		changes = append(changes, SAMLRoleChange{
			Field: "bound_subjects",
			From:  aSubjects,
			To:    bSubjects,
		})
	}

	aPolicies := strings.Join(a.TokenPolicies, ",")
	bPolicies := strings.Join(b.TokenPolicies, ",")
	if aPolicies != bPolicies {
		changes = append(changes, SAMLRoleChange{
			Field: "token_policies",
			From:  aPolicies,
			To:    bPolicies,
		})
	}

	return changes
}

// PrintSAMLDiff writes the SAML role diff to stdout.
func PrintSAMLDiff(a, b *vault.SAMLRoleInfo) {
	FprintSAMLDiff(os.Stdout, a, b)
}

// FprintSAMLDiff writes the SAML role diff to the given writer.
func FprintSAMLDiff(w io.Writer, a, b *vault.SAMLRoleInfo) {
	changes := CompareSAMLRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No SAML role changes detected.")
		return
	}
	fmt.Fprintln(w, "SAML Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
