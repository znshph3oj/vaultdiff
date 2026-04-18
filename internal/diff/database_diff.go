package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

type DatabaseRoleChange struct {
	Field    string
	OldValue string
	NewValue string
}

func CompareDatabaseRoles(a, b *vault.DatabaseRoleInfo) []DatabaseRoleChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []DatabaseRoleChange
	if a.DBName != b.DBName {
		changes = append(changes, DatabaseRoleChange{Field: "db_name", OldValue: a.DBName, NewValue: b.DBName})
	}
	if a.DefaultTTL != b.DefaultTTL {
		changes = append(changes, DatabaseRoleChange{Field: "default_ttl", OldValue: fmt.Sprintf("%d", a.DefaultTTL), NewValue: fmt.Sprintf("%d", b.DefaultTTL)})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, DatabaseRoleChange{Field: "max_ttl", OldValue: fmt.Sprintf("%d", a.MaxTTL), NewValue: fmt.Sprintf("%d", b.MaxTTL)})
	}
	return changes
}

func PrintDatabaseDiff(changes []DatabaseRoleChange) {
	FprintDatabaseDiff(os.Stdout, changes)
}

func FprintDatabaseDiff(w io.Writer, changes []DatabaseRoleChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No database role changes detected.")
		return
	}
	fmt.Fprintln(w, "Database Role Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  [~] %s: %q -> %q\n", c.Field, c.OldValue, c.NewValue)
	}
}
