package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType represents the type of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level change between two secret versions.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff result between two secret versions.
type Result struct {
	Path       string
	FromVersion int
	ToVersion   int
	Changes    []Change
}

// HasChanges returns true if there are any non-unchanged entries.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a brief string describing the number of additions,
// removals, and modifications in the result.
func (r *Result) Summary() string {
	var added, removed, modified int
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("+%d -%d ~%d", added, removed, modified)
}

// Compare computes the diff between two secret data maps.
func Compare(path string, fromVersion, toVersion int, from, to map[string]interface{}) *Result {
	result := &Result{
		Path:        path,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
	}

	keys := unionKeys(from, to)
	sort.Strings(keys)

	for _, key := range keys {
		oldVal, inFrom := from[key]
		newVal, inTo := to[key]

		switch {
		case inFrom && !inTo:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: fmt.Sprintf("%v", oldVal),
			})
		case !inFrom && inTo:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				NewValue: fmt.Sprintf("%v", newVal),
			})
		default:
			oldStr := fmt.Sprintf("%v", oldVal)
			newStr := fmt.Sprintf("%v", newVal)
			changeType := Unchanged
			if oldStr != newStr {
				changeType = Modified
			}
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     changeType,
				OldValue: oldStr,
				NewValue: newStr,
			})
		}
	}

	return result
}

func unionKeys(a, b map[string]interface{}) []string {
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
	return keys
}

// MaskValue replaces a secret value with asterisks.
func MaskValue(v string) string {
	if len(v) == 0 {
		return ""
	}
	return strings.Repeat("*", 8)
}
