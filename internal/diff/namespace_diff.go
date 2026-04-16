package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// NamespaceDiff represents added/removed namespaces between two snapshots.
type NamespaceDiff struct {
	Added   []string
	Removed []string
}

// CompareNamespaces diffs two slices of NamespaceInfo by path.
func CompareNamespaces(before, after []vault.NamespaceInfo) NamespaceDiff {
	bSet := make(map[string]struct{}, len(before))
	aSet := make(map[string]struct{}, len(after))
	for _, n := range before {
		bSet[n.Path] = struct{}{}
	}
	for _, n := range after {
		aSet[n.Path] = struct{}{}
	}

	var d NamespaceDiff
	for _, n := range after {
		if _, ok := bSet[n.Path]; !ok {
			d.Added = append(d.Added, n.Path)
		}
	}
	for _, n := range before {
		if _, ok := aSet[n.Path]; !ok {
			d.Removed = append(d.Removed, n.Path)
		}
	}
	return d
}

// PrintNamespaceDiff writes the diff to stdout.
func PrintNamespaceDiff(d NamespaceDiff) {
	FprintNamespaceDiff(os.Stdout, d)
}

// FprintNamespaceDiff writes the diff to the given writer.
func FprintNamespaceDiff(w io.Writer, d NamespaceDiff) {
	if len(d.Added) == 0 && len(d.Removed) == 0 {
		fmt.Fprintln(w, "No namespace changes detected.")
		return
	}
	for _, p := range d.Added {
		fmt.Fprintf(w, "+ %s\n", p)
	}
	for _, p := range d.Removed {
		fmt.Fprintf(w, "- %s\n", p)
	}
}
