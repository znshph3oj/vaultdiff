package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// SysConfigChange represents a changed field in system config.
type SysConfigChange struct {
	Field string
	From  string
	To    string
}

// CompareSysConfig returns changes between two SysConfig snapshots.
func CompareSysConfig(a, b *vault.SysConfig) []SysConfigChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []SysConfigChange
	if a.DefaultLeaseTTL != b.DefaultLeaseTTL {
		changes = append(changes, SysConfigChange{Field: "default_lease_ttl", From: a.DefaultLeaseTTL, To: b.DefaultLeaseTTL})
	}
	if a.MaxLeaseTTL != b.MaxLeaseTTL {
		changes = append(changes, SysConfigChange{Field: "max_lease_ttl", From: a.MaxLeaseTTL, To: b.MaxLeaseTTL})
	}
	if a.ForceNoCache != b.ForceNoCache {
		changes = append(changes, SysConfigChange{
			Field: "force_no_cache",
			From:  fmt.Sprintf("%v", a.ForceNoCache),
			To:    fmt.Sprintf("%v", b.ForceNoCache),
		})
	}
	return changes
}

// PrintSysConfigDiff prints sys config diff to stdout.
func PrintSysConfigDiff(changes []SysConfigChange) {
	FprintSysConfigDiff(os.Stdout, changes)
}

// FprintSysConfigDiff writes sys config diff to the given writer.
func FprintSysConfigDiff(w io.Writer, changes []SysConfigChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No sys config changes.")
		return
	}
	fmt.Fprintln(w, "Sys Config Changes:")
	for _, c := range changes {
		fmt.Fprintf(w, "  [~] %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
