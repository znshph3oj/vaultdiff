package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// TerraformChange represents a single field change between two TerraformRoleInfo values.
type TerraformChange struct {
	Field string
	From  string
	To    string
}

// CompareTerraformRoles returns a list of changes between two Terraform role configurations.
// Returns nil if both inputs are nil.
func CompareTerraformRoles(a, b *vault.TerraformRoleInfo) []TerraformChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.TerraformRoleInfo{}
	}
	if b == nil {
		b = &vault.TerraformRoleInfo{}
	}

	var changes []TerraformChange
	check := func(field, from, to string) {
		if from != to {
			changes = append(changes, TerraformChange{Field: field, From: from, To: to})
		}
	}

	check("organization", a.Organization, b.Organization)
	check("team_id", a.TeamID, b.TeamID)
	check("ttl", a.TTL, b.TTL)
	check("max_ttl", a.MaxTTL, b.MaxTTL)
	check("token_type", a.TokenType, b.TokenType)

	return changes
}

// PrintTerraformDiff writes a Terraform role diff to stdout.
func PrintTerraformDiff(role string, changes []TerraformChange) {
	FprintTerraformDiff(os.Stdout, role, changes)
}

// FprintTerraformDiff writes a Terraform role diff to the given writer.
func FprintTerraformDiff(w io.Writer, role string, changes []TerraformChange) {
	if len(changes) == 0 {
		fmt.Fprintf(w, "terraform role %q: no changes\n", role)
		return
	}
	fmt.Fprintf(w, "terraform role %q changes:\n", role)
	for _, c := range changes {
		fmt.Fprintf(w, "  %-14s %q -> %q\n", c.Field+":", c.From, c.To)
	}
}
