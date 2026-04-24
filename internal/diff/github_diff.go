package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// GitHubRoleChange describes a single field change in a GitHub auth role.
type GitHubRoleChange struct {
	Field string
	From  string
	To    string
}

// CompareGitHubRoles returns the list of changes between two GitHubRoleInfo values.
func CompareGitHubRoles(a, b *vault.GitHubRoleInfo) []GitHubRoleChange {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = &vault.GitHubRoleInfo{}
	}
	if b == nil {
		b = &vault.GitHubRoleInfo{}
	}

	var changes []GitHubRoleChange

	if a.TTL != b.TTL {
		changes = append(changes, GitHubRoleChange{Field: "ttl", From: a.TTL, To: b.TTL})
	}
	if a.MaxTTL != b.MaxTTL {
		changes = append(changes, GitHubRoleChange{Field: "max_ttl", From: a.MaxTTL, To: b.MaxTTL})
	}

	aSet := sliceToStrSet(a.Policies)
	bSet := sliceToStrSet(b.Policies)
	for p := range bSet {
		if !aSet[p] {
			changes = append(changes, GitHubRoleChange{Field: "policy", From: "", To: p})
		}
	}
	for p := range aSet {
		if !bSet[p] {
			changes = append(changes, GitHubRoleChange{Field: "policy", From: p, To: ""})
		}
	}

	return changes
}

// PrintGitHubDiff writes the diff to stdout.
func PrintGitHubDiff(team string, changes []GitHubRoleChange) {
	FprintGitHubDiff(os.Stdout, team, changes)
}

// FprintGitHubDiff writes the diff to the provided writer.
func FprintGitHubDiff(w io.Writer, team string, changes []GitHubRoleChange) {
	if len(changes) == 0 {
		fmt.Fprintf(w, "github role %q: no changes\n", team)
		return
	}
	fmt.Fprintf(w, "github role %q changes:\n", team)
	for _, c := range changes {
		switch {
		case c.From == "":
			fmt.Fprintf(w, "  + %s: %s\n", c.Field, c.To)
		case c.To == "":
			fmt.Fprintf(w, "  - %s: %s\n", c.Field, c.From)
		default:
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Field, c.From, c.To)
		}
	}
	_ = strings.TrimSpace("") // satisfy import
}
