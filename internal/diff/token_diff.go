package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// TokenDiff represents the comparison between two TokenInfo snapshots.
type TokenDiff struct {
	PoliciesAdded   []string
	PoliciesRemoved []string
	TTLChanged      bool
	OldTTL          int
	NewTTL          int
	RenewableChanged bool
	OldRenewable    bool
	NewRenewable    bool
}

// CompareTokens computes the diff between two TokenInfo values.
func CompareTokens(a, b *vault.TokenInfo) *TokenDiff {
	if a == nil {
		a = &vault.TokenInfo{}
	}
	if b == nil {
		b = &vault.TokenInfo{}
	}

	d := &TokenDiff{}

	aSet := make(map[string]bool)
	for _, p := range a.Policies {
		aSet[p] = true
	}
	bSet := make(map[string]bool)
	for _, p := range b.Policies {
		bSet[p] = true
	}
	for p := range bSet {
		if !aSet[p] {
			d.PoliciesAdded = append(d.PoliciesAdded, p)
		}
	}
	for p := range aSet {
		if !bSet[p] {
			d.PoliciesRemoved = append(d.PoliciesRemoved, p)
		}
	}

	if a.TTL != b.TTL {
		d.TTLChanged = true
		d.OldTTL = a.TTL
		d.NewTTL = b.TTL
	}

	if a.Renewable != b.Renewable {
		d.RenewableChanged = true
		d.OldRenewable = a.Renewable
		d.NewRenewable = b.Renewable
	}

	return d
}

// PrintTokenDiff writes a human-readable token diff to stdout.
func PrintTokenDiff(d *TokenDiff) {
	FprintTokenDiff(os.Stdout, d)
}

// FprintTokenDiff writes a human-readable token diff to the given writer.
func FprintTokenDiff(w io.Writer, d *TokenDiff) {
	fmt.Fprintln(w, "=== Token Diff ===")
	if len(d.PoliciesAdded) > 0 {
		fmt.Fprintf(w, "+ policies added:   %s\n", strings.Join(d.PoliciesAdded, ", "))
	}
	if len(d.PoliciesRemoved) > 0 {
		fmt.Fprintf(w, "- policies removed: %s\n", strings.Join(d.PoliciesRemoved, ", "))
	}
	if d.TTLChanged {
		fmt.Fprintf(w, "~ ttl: %d -> %d\n", d.OldTTL, d.NewTTL)
	}
	if d.RenewableChanged {
		fmt.Fprintf(w, "~ renewable: %v -> %v\n", d.OldRenewable, d.NewRenewable)
	}
	if len(d.PoliciesAdded) == 0 && len(d.PoliciesRemoved) == 0 && !d.TTLChanged && !d.RenewableChanged {
		fmt.Fprintln(w, "  (no changes)")
	}
}
