package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

type SSHRoleChange struct {
	Field string
	From  string
	To    string
}

func CompareSSHRoles(a, b *vault.SSHRoleInfo) []SSHRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []SSHRoleChange
	check := func(field, from, to string) {
		if from != to {
			changes = append(changes, SSHRoleChange{Field: field, From: from, To: to})
		}
	}
	check("key_type", a.KeyType, b.KeyType)
	check("default_user", a.DefaultUser, b.DefaultUser)
	check("allowed_users", a.AllowedUsers, b.AllowedUsers)
	check("ttl", a.TTL, b.TTL)
	check("max_ttl", a.MaxTTL, b.MaxTTL)
	check("allowed_domains", a.AllowedDomains, b.AllowedDomains)
	check("cidr_list", a.CIDRList, b.CIDRList)
	check("allowed_extensions", a.AllowedExtensions, b.AllowedExtensions)

	aP := strings.Join(a.Policies, ",")
	bP := strings.Join(b.Policies, ",")
	check("policies", aP, bP)

	return changes
}

func PrintSSHDiff(changes []SSHRoleChange) {
	FprintSSHDiff(os.Stdout, changes)
}

func FprintSSHDiff(w io.Writer, changes []SSHRoleChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No SSH role changes detected.")
		return
	}
	fmt.Fprintln(w, "SSH Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  [~] %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
