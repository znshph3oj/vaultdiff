package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/vaultdiff/internal/vault"
)

type PolicyChange struct {
	Name   string
	Status string // added, removed, modified, unchanged
	OldRules string
	NewRules string
}

func ComparePolicies(oldPolicies, newPolicies map[string]*vault.PolicyInfo) []PolicyChange {
	var changes []PolicyChange
	keys := map[string]struct{}{}
	for k := range oldPolicies {
		keys[k] = struct{}{}
	}
	for k := range newPolicies {
		keys[k] = struct{}{}
	}
	for k := range keys {
		old, inOld := oldPolicies[k]
		new, inNew := newPolicies[k]
		switch {
		case inOld && !inNew:
			changes = append(changes, PolicyChange{Name: k, Status: "removed", OldRules: old.Rules})
		case !inOld && inNew:
			changes = append(changes, PolicyChange{Name: k, Status: "added", NewRules: new.Rules})
		case old.Rules != new.Rules:
			changes = append(changes, PolicyChange{Name: k, Status: "modified", OldRules: old.Rules, NewRules: new.Rules})
		default:
			changes = append(changes, PolicyChange{Name: k, Status: "unchanged"})
		}
	}
	return changes
}

func PrintPolicyDiff(changes []PolicyChange) {
	FprintPolicyDiff(os.Stdout, changes)
}

func FprintPolicyDiff(w io.Writer, changes []PolicyChange) {
	fmt.Fprintln(w, "Policy Diff:")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range changes {
		switch c.Status {
		case "added":
			fmt.Fprintf(w, "[+] %s (added)\n", c.Name)
		case "removed":
			fmt.Fprintf(w, "[-] %s (removed)\n", c.Name)
		case "modified":
			fmt.Fprintf(w, "[~] %s (modified)\n", c.Name)
			fmt.Fprintf(w, "    old: %s\n", c.OldRules)
			fmt.Fprintf(w, "    new: %s\n", c.NewRules)
		}
	}
}
