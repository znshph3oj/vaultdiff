package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// JWTRoleChange describes a single field change in a JWT role.
type JWTRoleChange struct {
	Field string
	From  string
	To    string
}

// CompareJWTRoles returns a list of changes between two JWTRoleInfo structs.
func CompareJWTRoles(a, b *vault.JWTRoleInfo) []JWTRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []JWTRoleChange

	check := func(field, from, to string) {
		if from != to {
			changes = append(changes, JWTRoleChange{Field: field, From: from, To: to})
		}
	}

	check("role_type", a.RoleType, b.RoleType)
	check("user_claim", a.UserClaim, b.UserClaim)
	check("groups_claim", a.GroupsClaim, b.GroupsClaim)
	check("bound_subject", a.BoundSubject, b.BoundSubject)
	check("ttl", a.TTL, b.TTL)
	check("max_ttl", a.MaxTTL, b.MaxTTL)
	check("bound_audiences",
		strings.Join(a.BoundAudiences, ","),
		strings.Join(b.BoundAudiences, ","))
	check("token_policies",
		strings.Join(a.TokenPolicies, ","),
		strings.Join(b.TokenPolicies, ","))

	return changes
}

// PrintJWTDiff writes a JWT role diff to stdout.
func PrintJWTDiff(a, b *vault.JWTRoleInfo) {
	FprintJWTDiff(os.Stdout, a, b)
}

// FprintJWTDiff writes a JWT role diff to the given writer.
func FprintJWTDiff(w io.Writer, a, b *vault.JWTRoleInfo) {
	changes := CompareJWTRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No JWT role changes detected.")
		return
	}
	fmt.Fprintln(w, "JWT Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  %-20s  - %s\n", c.Field, c.From)
		fmt.Fprintf(w, "  %-20s  + %s\n", "", c.To)
	}
}
