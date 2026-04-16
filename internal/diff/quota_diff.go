package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultdiff/internal/vault"
)

// QuotaChange describes a single changed field between two QuotaInfo values.
type QuotaChange struct {
	Field string
	From  string
	To    string
}

// CompareQuotas returns the list of changes between two QuotaInfo snapshots.
func CompareQuotas(a, b *vault.QuotaInfo) []QuotaChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []QuotaChange
	if a.Type != b.Type {
		changes = append(changes, QuotaChange{Field: "type", From: a.Type, To: b.Type})
	}
	if a.Rate != b.Rate {
		changes = append(changes, QuotaChange{
			Field: "rate",
			From:  fmt.Sprintf("%.2f", a.Rate),
			To:    fmt.Sprintf("%.2f", b.Rate),
		})
	}
	if a.Interval != b.Interval {
		changes = append(changes, QuotaChange{
			Field: "interval",
			From:  fmt.Sprintf("%.2f", a.Interval),
			To:    fmt.Sprintf("%.2f", b.Interval),
		})
	}
	if a.BlockInterval != b.BlockInterval {
		changes = append(changes, QuotaChange{
			Field: "block_interval",
			From:  fmt.Sprintf("%.2f", a.BlockInterval),
			To:    fmt.Sprintf("%.2f", b.BlockInterval),
		})
	}
	return changes
}

// PrintQuotaDiff writes quota changes to stdout.
func PrintQuotaDiff(name string, changes []QuotaChange) {
	FprintQuotaDiff(os.Stdout, name, changes)
}

// FprintQuotaDiff writes quota changes to the given writer.
func FprintQuotaDiff(w io.Writer, name string, changes []QuotaChange) {
	if len(changes) == 0 {
		fmt.Fprintf(w, "quota %q: no changes\n", name)
		return
	}
	fmt.Fprintf(w, "quota %q changes:\n", name)
	for _, c := range changes {
		fmt.Fprintf(w, "  %s: %s -> %s\n", c.Field, c.From, c.To)
	}
}
