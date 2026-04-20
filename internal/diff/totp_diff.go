package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/your/vaultdiff/internal/vault"
)

// TOTPChange represents a single field change between two TOTP key configs.
type TOTPChange struct {
	Field string
	From  interface{}
	To    interface{}
}

// CompareTOTPKeys compares two TOTPKeyInfo structs and returns a list of changes.
func CompareTOTPKeys(a, b *vault.TOTPKeyInfo) []TOTPChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []TOTPChange
	if a.AccountName != b.AccountName {
		changes = append(changes, TOTPChange{"account_name", a.AccountName, b.AccountName})
	}
	if a.Algorithm != b.Algorithm {
		changes = append(changes, TOTPChange{"algorithm", a.Algorithm, b.Algorithm})
	}
	if a.Digits != b.Digits {
		changes = append(changes, TOTPChange{"digits", a.Digits, b.Digits})
	}
	if a.Issuer != b.Issuer {
		changes = append(changes, TOTPChange{"issuer", a.Issuer, b.Issuer})
	}
	if a.Period != b.Period {
		changes = append(changes, TOTPChange{"period", a.Period, b.Period})
	}
	if a.QRSize != b.QRSize {
		changes = append(changes, TOTPChange{"qr_size", a.QRSize, b.QRSize})
	}
	return changes
}

// PrintTOTPDiff writes the TOTP key diff to stdout.
func PrintTOTPDiff(changes []TOTPChange) {
	FprintTOTPDiff(os.Stdout, changes)
}

// FprintTOTPDiff writes the TOTP key diff to the given writer.
func FprintTOTPDiff(w io.Writer, changes []TOTPChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No TOTP key changes detected.")
		return
	}
	fmt.Fprintln(w, "TOTP Key Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %v -> %v\n", c.Field, c.From, c.To)
	}
}
