package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

type LDAPRoleChange struct {
	Field string
	From  string
	To    string
}

func CompareLDAPRoles(a, b *vault.LDAPRoleInfo) []LDAPRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []LDAPRoleChange
	check := func(field, from, to string) {
		if from != to {
			changes = append(changes, LDAPRoleChange{Field: field, From: from, To: to})
		}
	}
	check("group_dn", a.GroupDN, b.GroupDN)
	check("user_dn", a.UserDN, b.UserDN)
	check("bind_dn", a.BindDN, b.BindDN)
	check("group_filter", a.GroupFilter, b.GroupFilter)
	check("group_attr", a.GroupAttr, b.GroupAttr)
	check("user_attr", a.UserAttr, b.UserAttr)
	check("ttl", a.TTL, b.TTL)
	check("max_ttl", a.MaxTTL, b.MaxTTL)
	check("policies", strings.Join(a.Policies, ","), strings.Join(b.Policies, ","))
	return changes
}

func PrintLDAPDiff(a, b *vault.LDAPRoleInfo) {
	FprintLDAPDiff(os.Stdout, a, b)
}

func FprintLDAPDiff(w io.Writer, a, b *vault.LDAPRoleInfo) {
	changes := CompareLDAPRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes in LDAP role.")
		return
	}
	fmt.Fprintln(w, "LDAP Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
