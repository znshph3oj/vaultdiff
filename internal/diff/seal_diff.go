package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// SealDiffEntry represents a changed field between two seal statuses.
type SealDiffEntry struct {
	Field string
	From  string
	To    string
}

// CompareSealStatus diffs two SealStatus values and returns changed fields.
func CompareSealStatus(a, b *vault.SealStatus) []SealDiffEntry {
	if a == nil || b == nil {
		return nil
	}
	var entries []SealDiffEntry

	if a.Sealed != b.Sealed {
		entries = append(entries, SealDiffEntry{"sealed", fmt.Sprintf("%v", a.Sealed), fmt.Sprintf("%v", b.Sealed)})
	}
	if a.Initialized != b.Initialized {
		entries = append(entries, SealDiffEntry{"initialized", fmt.Sprintf("%v", a.Initialized), fmt.Sprintf("%v", b.Initialized)})
	}
	if a.Version != b.Version {
		entries = append(entries, SealDiffEntry{"version", a.Version, b.Version})
	}
	if a.ClusterName != b.ClusterName {
		entries = append(entries, SealDiffEntry{"cluster_name", a.ClusterName, b.ClusterName})
	}
	if a.T != b.T {
		entries = append(entries, SealDiffEntry{"threshold", fmt.Sprintf("%d", a.T), fmt.Sprintf("%d", b.T)})
	}
	if a.N != b.N {
		entries = append(entries, SealDiffEntry{"shares", fmt.Sprintf("%d", a.N), fmt.Sprintf("%d", b.N)})
	}
	return entries
}

// PrintSealDiff writes a seal status diff to stdout.
func PrintSealDiff(a, b *vault.SealStatus) {
	FprintSealDiff(os.Stdout, a, b)
}

// FprintSealDiff writes a seal status diff to the given writer.
func FprintSealDiff(w io.Writer, a, b *vault.SealStatus) {
	entries := CompareSealStatus(a, b)
	if len(entries) == 0 {
		fmt.Fprintln(w, "seal status: no changes")
		return
	}
	fmt.Fprintln(w, "seal status diff:")
	for _, e := range entries {
		fmt.Fprintf(w, "  %-16s %s -> %s\n", e.Field, e.From, e.To)
	}
}
