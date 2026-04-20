package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// ConsulRoleChange represents a single field-level change between two Consul roles.
type ConsulRoleChange struct {
	Field  string
	OldVal string
	NewVal string
}

// CompareConsulRoles returns a list of field-level changes between two ConsulRoleInfo structs.
func CompareConsulRoles(a, b *vault.ConsulRoleInfo) []ConsulRoleChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.ConsulRoleInfo{}
	}
	if b == nil {
		b = &vault.ConsulRoleInfo{}
	}

	var changes []ConsulRoleChange

	if a.TokenType != b.TokenType {
		changes = append(changes, ConsulRoleChange{Field: "token_type", OldVal: a.TokenType, NewVal: b.TokenType})
	}
	if a.TTL != b.TTL {
		changes = append(changes, ConsulRoleChange{Field: "ttl", OldVal: a.TTL, NewVal: b.TTL})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, ConsulRoleChange{Field: "max_ttl", OldVal: a.MaxTTL, NewVal: b.MaxTTL})
	}
	if fmt.Sprintf("%v", a.Local) != fmt.Sprintf("%v", b.Local) {
		changes = append(changes, ConsulRoleChange{
			Field:  "local",
			OldVal: fmt.Sprintf("%v", a.Local),
			NewVal: fmt.Sprintf("%v", b.Local),
		})
	}

	aSet := sliceToStrSet(a.Policies)
	bSet := sliceToStrSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, ConsulRoleChange{Field: "policy:+" + p, OldVal: "", NewVal: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, ConsulRoleChange{Field: "policy:-" + p, OldVal: p, NewVal: ""})
		}
	}

	return changes
}

// PrintConsulDiff writes a human-readable diff of two Consul roles to stdout.
func PrintConsulDiff(a, b *vault.ConsulRoleInfo) {
	FprintConsulDiff(os.Stdout, a, b)
}

// FprintConsulDiff writes a human-readable diff of two Consul roles to w.
func FprintConsulDiff(w io.Writer, a, b *vault.ConsulRoleInfo) {
	changes := CompareConsulRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes in Consul role.")
		return
	}
	fmt.Fprintln(w, "Consul Role Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		if c.OldVal == "" {
			fmt.Fprintf(w, "  + %-20s %s\n", c.Field, c.NewVal)
		} else if c.NewVal == "" {
			fmt.Fprintf(w, "  - %-20s %s\n", c.Field, c.OldVal)
		} else {
			fmt.Fprintf(w, "  ~ %-20s %s -> %s\n", c.Field, c.OldVal, c.NewVal)
		}
	}
}
