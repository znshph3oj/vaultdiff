package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// OIDCRoleChange describes a single changed field between two OIDC role configs.
type OIDCRoleChange struct {
	Field string
	From  string
	To    string
}

// CompareOIDCRoles compares two OIDCRoleInfo structs and returns a list of changes.
func CompareOIDCRoles(a, b *vault.OIDCRoleInfo) []OIDCRoleChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.OIDCRoleInfo{}
	}
	if b == nil {
		b = &vault.OIDCRoleInfo{}
	}

	var changes []OIDCRoleChange

	if a.UserClaim != b.UserClaim {
		changes = append(changes, OIDCRoleChange{Field: "user_claim", From: a.UserClaim, To: b.UserClaim})
	}
	if a.TokenTTL != b.TokenTTL {
		changes = append(changes, OIDCRoleChange{
			Field: "token_ttl",
			From:  fmt.Sprintf("%d", a.TokenTTL),
			To:    fmt.Sprintf("%d", b.TokenTTL),
		})
	}
	if a.TokenMaxTTL != b.TokenMaxTTL {
		changes = append(changes, OIDCRoleChange{
			Field: "token_max_ttl",
			From:  fmt.Sprintf("%d", a.TokenMaxTTL),
			To:    fmt.Sprintf("%d", b.TokenMaxTTL),
		})
	}

	aAud := strings.Join(a.BoundAudiences, ",")
	bAud := strings.Join(b.BoundAudiences, ",")
	if aAud != bAud {
		changes = append(changes, OIDCRoleChange{Field: "bound_audiences", From: aAud, To: bAud})
	}

	aRedir := strings.Join(a.AllowedRedirects, ",")
	bRedir := strings.Join(b.AllowedRedirects, ",")
	if aRedir != bRedir {
		changes = append(changes, OIDCRoleChange{Field: "allowed_redirect_uris", From: aRedir, To: bRedir})
	}

	aPolicy := strings.Join(a.TokenPolicies, ",")
	bPolicy := strings.Join(b.TokenPolicies, ",")
	if aPolicy != bPolicy {
		changes = append(changes, OIDCRoleChange{Field: "token_policies", From: aPolicy, To: bPolicy})
	}

	return changes
}

// PrintOIDCDiff writes a human-readable diff of OIDC role changes to stdout.
func PrintOIDCDiff(changes []OIDCRoleChange) {
	FprintOIDCDiff(os.Stdout, changes)
}

// FprintOIDCDiff writes a human-readable diff of OIDC role changes to w.
func FprintOIDCDiff(w io.Writer, changes []OIDCRoleChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No OIDC role changes detected.")
		return
	}
	fmt.Fprintln(w, "OIDC Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
