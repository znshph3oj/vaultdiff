package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// GCPKMSChange represents a single changed field in a GCP KMS key.
type GCPKMSChange struct {
	Field string
	From  string
	To    string
}

// CompareGCPKMSKeys compares two GCPKMSKeyInfo structs and returns the list of changes.
func CompareGCPKMSKeys(a, b *vault.GCPKMSKeyInfo) []GCPKMSChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []GCPKMSChange
	check := func(field, from, to string) {
		if from != to {
			changes = append(changes, GCPKMSChange{Field: field, From: from, To: to})
		}
	}
	check("algorithm", a.Algorithm, b.Algorithm)
	check("protection_level", a.ProtectionLevel, b.ProtectionLevel)
	check("rotation_period", a.RotationPeriod, b.RotationPeriod)
	check("key_ring", a.KeyRing, b.KeyRing)
	check("crypto_key", a.CryptoKey, b.CryptoKey)
	if a.MinVersion != b.MinVersion {
		changes = append(changes, GCPKMSChange{
			Field: "min_version",
			From:  fmt.Sprintf("%d", a.MinVersion),
			To:    fmt.Sprintf("%d", b.MinVersion),
		})
	}
	return changes
}

// PrintGCPKMSDiff prints a diff of two GCP KMS keys to stdout.
func PrintGCPKMSDiff(a, b *vault.GCPKMSKeyInfo) {
	FprintGCPKMSDiff(os.Stdout, a, b)
}

// FprintGCPKMSDiff writes a diff of two GCP KMS keys to the given writer.
func FprintGCPKMSDiff(w io.Writer, a, b *vault.GCPKMSKeyInfo) {
	changes := CompareGCPKMSKeys(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "No GCP KMS key changes detected.")
		return
	}
	fmt.Fprintln(w, "GCP KMS Key Diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  ~ %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
