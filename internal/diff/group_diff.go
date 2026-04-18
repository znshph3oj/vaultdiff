package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// GroupChange represents a single field-level change in an identity group.
type GroupChange struct {
	Field string
	From  string
	To    string
}

// CompareGroups returns the list of changes between two identity groups.
func CompareGroups(a, b *vault.IdentityGroup) []GroupChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []GroupChange

	if a.Name != b.Name {
		changes = append(changes, GroupChange{Field: "name", From: a.Name, To: b.Name})
	}
	if a.Type != b.Type {
		changes = append(changes, GroupChange{Field: "type", From: a.Type, To: b.Type})
	}

	aPolSet := sliceToSet(a.Policies)
	bPolSet := sliceToSet(b.Policies)
	for p := range bPolSet {
		if !aPolSet[p] {
			changes = append(changes, GroupChange{Field: "policy", From: "", To: p})
		}
	}
	for p := range aPolSet {
		if !bPolSet[p] {
			changes = append(changes, GroupChange{Field: "policy", From: p, To: ""})
		}
	}

	aMembers := sliceToSet(a.MemberEntityIDs)
	bMembers := sliceToSet(b.MemberEntityIDs)
	for m := range bMembers {
		if !aMembers[m] {
			changes = append(changes, GroupChange{Field: "member_entity_id", From: "", To: m})
		}
	}
	for m := range aMembers {
		if !bMembers[m] {
			changes = append(changes, GroupChange{Field: "member_entity_id", From: m, To: ""})
		}
	}

	return changes
}

// PrintGroupDiff prints group changes to stdout.
func PrintGroupDiff(changes []GroupChange) {
	FprintGroupDiff(os.Stdout, changes)
}

// FprintGroupDiff writes group changes to the provided writer.
func FprintGroupDiff(w io.Writer, changes []GroupChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No group changes detected.")
		return
	}
	fmt.Fprintln(w, "Group Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch {
		case c.From == "":
			fmt.Fprintf(w, "  + [%s] %s\n", c.Field, c.To)
		case c.To == "":
			fmt.Fprintf(w, "  - [%s] %s\n", c.Field, c.From)
		default:
			fmt.Fprintf(w, "  ~ [%s] %s -> %s\n", c.Field, c.From, c.To)
		}
	}
}
