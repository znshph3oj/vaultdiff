package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultdiff/internal/vault"
)

type AWSDiffEntry struct {
	Field    string
	OldValue string
	NewValue string
}

func CompareAWSRoles(a, b *vault.AWSRoleInfo) []AWSDiffEntry {
	if a == nil && b == nil {
		return nil
	}
	var changes []AWSDiffEntry
	if a == nil {
		a = &vault.AWSRoleInfo{}
	}
	if b == nil {
		b = &vault.AWSRoleInfo{}
	}
	if a.CredentialType != b.CredentialType {
		changes = append(changes, AWSDiffEntry{"credential_type", a.CredentialType, b.CredentialType})
	}
	if a.RoleARNs != b.RoleARNs {
		changes = append(changes, AWSDiffEntry{"role_arns", a.RoleARNs, b.RoleARNs})
	}
	if a.PolicyARNs != b.PolicyARNs {
		changes = append(changes, AWSDiffEntry{"policy_arns", a.PolicyARNs, b.PolicyARNs})
	}
	if a.DefaultSTSTTL != b.DefaultSTSTTL {
		changes = append(changes, AWSDiffEntry{"default_sts_ttl", fmt.Sprintf("%d", a.DefaultSTSTTL), fmt.Sprintf("%d", b.DefaultSTSTTL)})
	}
	if a.MaxSTSTTL != b.MaxSTSTTL {
		changes = append(changes, AWSDiffEntry{"max_sts_ttl", fmt.Sprintf("%d", a.MaxSTSTTL), fmt.Sprintf("%d", b.MaxSTSTTL)})
	}
	return changes
}

func PrintAWSDiff(changes []AWSDiffEntry) {
	FprintAWSDiff(os.Stdout, changes)
}

func FprintAWSDiff(w io.Writer, changes []AWSDiffEntry) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No AWS role changes detected.")
		return
	}
	fmt.Fprintln(w, "AWS Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.OldValue, c.NewValue)
	}
}
