package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// KubernetesRoleChange describes a single field-level change in a Kubernetes role binding.
type KubernetesRoleChange struct {
	Field string
	Old   string
	New   string
}

// CompareKubernetesRoles returns the list of changes between two KubernetesRoleBinding values.
func CompareKubernetesRoles(a, b *vault.KubernetesRoleBinding) []KubernetesRoleChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.KubernetesRoleBinding{}
	}
	if b == nil {
		b = &vault.KubernetesRoleBinding{}
	}

	var changes []KubernetesRoleChange

	cmpStr := func(field, old, new string) {
		if old != new {
			changes = append(changes, KubernetesRoleChange{Field: field, Old: old, New: new})
		}
	}
	cmpSlice := func(field string, old, new []string) {
		if strings.Join(old, ",") != strings.Join(new, ",") {
			changes = append(changes, KubernetesRoleChange{
				Field: field,
				Old:   strings.Join(old, ", "),
				New:   strings.Join(new, ", "),
			})
		}
	}

	cmpStr("ttl", a.TTL, b.TTL)
	cmpStr("max_ttl", a.MaxTTL, b.MaxTTL)
	cmpSlice("bound_service_account_names", a.BoundServiceAccountNames, b.BoundServiceAccountNames)
	cmpSlice("bound_service_account_namespaces", a.BoundServiceAccountNamespaces, b.BoundServiceAccountNamespaces)
	cmpSlice("token_policies", a.Policies, b.Policies)

	return changes
}

// PrintKubernetesRoleDiff prints the diff to stdout.
func PrintKubernetesRoleDiff(a, b *vault.KubernetesRoleBinding) {
	FprintKubernetesRoleDiff(os.Stdout, a, b)
}

// FprintKubernetesRoleDiff writes the diff to the given writer.
func FprintKubernetesRoleDiff(w io.Writer, a, b *vault.KubernetesRoleBinding) {
	changes := CompareKubernetesRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes in Kubernetes role binding.")
		return
	}
	fmt.Fprintln(w, "Kubernetes Role Binding Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  %-40s %s -> %s\n", c.Field+":", c.Old, c.New)
	}
}
