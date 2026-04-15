package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// TagChange represents a single change in a secret's custom metadata tag.
type TagChange struct {
	Key    string
	OldVal string
	NewVal string
	Type   ChangeType
}

// CompareTags diffs two SecretTags maps and returns the list of tag changes.
func CompareTags(before, after map[string]string) []TagChange {
	keys := unionStringKeys(before, after)
	var changes []TagChange
	for _, k := range keys {
		oldV, inOld := before[k]
		newV, inNew := after[k]
		switch {
		case inOld && !inNew:
			changes = append(changes, TagChange{Key: k, OldVal: oldV, Type: Removed})
		case !inOld && inNew:
			changes = append(changes, TagChange{Key: k, NewVal: newV, Type: Added})
		case oldV != newV:
			changes = append(changes, TagChange{Key: k, OldVal: oldV, NewVal: newV, Type: Modified})
		}
	}
	return changes
}

// PrintTagDiff writes a human-readable tag diff to stdout.
func PrintTagDiff(changes []TagChange) {
	FprintTagDiff(os.Stdout, changes)
}

// FprintTagDiff writes a human-readable tag diff to the given writer.
func FprintTagDiff(w io.Writer, changes []TagChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "  (no tag changes)")
		return
	}
	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + [tag] %s = %q\n", c.Key, c.NewVal)
		case Removed:
			fmt.Fprintf(w, "  - [tag] %s = %q\n", c.Key, c.OldVal)
		case Modified:
			fmt.Fprintf(w, "  ~ [tag] %s: %q -> %q\n", c.Key, c.OldVal, c.NewVal)
		}
	}
}

// unionStringKeys returns a sorted union of keys from two string maps.
func unionStringKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
