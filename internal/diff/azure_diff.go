package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

type AzureRoleChange struct {
	Field string
	From  string
	To    string
}

func CompareAzureRoles(a, b *vault.AzureRoleInfo) []AzureRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []AzureRoleChange
	if a.ApplicationObjectID != b.ApplicationObjectID {
		changes = append(changes, AzureRoleChange{"application_object_id", a.ApplicationObjectID, b.ApplicationObjectID})
	}
	if a.ClientID != b.ClientID {
		changes = append(changes, AzureRoleChange{"client_id", a.ClientID, b.ClientID})
	}
	if a.TTL != b.TTL {
		changes = append(changes, AzureRoleChange{"ttl", a.TTL, b.TTL})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, AzureRoleChange{"max_ttl", a.MaxTTL, b.MaxTTL})
	}
	if strings.Join(a.AzureRoles, ",") != strings.Join(b.AzureRoles, ",") {
		changes = append(changes, AzureRoleChange{"azure_roles", strings.Join(a.AzureRoles, ","), strings.Join(b.AzureRoles, ",")})
	}
	if strings.Join(a.AzureGroups, ",") != strings.Join(b.AzureGroups, ",") {
		changes = append(changes, AzureRoleChange{"azure_groups", strings.Join(a.AzureGroups, ","), strings.Join(b.AzureGroups, ",")})
	}
	return changes
}

func PrintAzureDiff(a, b *vault.AzureRoleInfo) {
	FprintAzureDiff(os.Stdout, a, b)
}

func FprintAzureDiff(w io.Writer, a, b *vault.AzureRoleInfo) {
	changes := CompareAzureRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes in Azure role.")
		return
	}
	fmt.Fprintln(w, "Azure Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
