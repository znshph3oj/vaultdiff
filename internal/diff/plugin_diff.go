package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// PluginChange represents a change in a plugin entry.
type PluginChange struct {
	Name   string
	Status string // "added", "removed"
}

// PluginDiffResult holds the result of comparing two plugin lists.
type PluginDiffResult struct {
	Changes []PluginChange
}

// ComparePlugins compares two sets of plugin names and returns added/removed entries.
func ComparePlugins(old, new *vault.PluginList) *PluginDiffResult {
	result := &PluginDiffResult{}

	if old == nil && new == nil {
		return result
	}

	oldSet := make(map[string]bool)
	newSet := make(map[string]bool)

	if old != nil {
		for _, name := range old.Plugins {
			oldSet[name] = true
		}
	}
	if new != nil {
		for _, name := range new.Plugins {
			newSet[name] = true
		}
	}

	// Detect removed
	for name := range oldSet {
		if !newSet[name] {
			result.Changes = append(result.Changes, PluginChange{Name: name, Status: "removed"})
		}
	}
	// Detect added
	for name := range newSet {
		if !oldSet[name] {
			result.Changes = append(result.Changes, PluginChange{Name: name, Status: "added"})
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		return result.Changes[i].Name < result.Changes[j].Name
	})

	return result
}

// PrintPluginDiff prints the plugin diff to stdout.
func PrintPluginDiff(result *PluginDiffResult) {
	FprintPluginDiff(os.Stdout, result)
}

// FprintPluginDiff writes the plugin diff to the given writer.
func FprintPluginDiff(w io.Writer, result *PluginDiffResult) {
	if result == nil || len(result.Changes) == 0 {
		fmt.Fprintln(w, "No plugin changes detected.")
		return
	}
	fmt.Fprintln(w, "Plugin Diff:")
	for _, c := range result.Changes {
		switch c.Status {
		case "added":
			fmt.Fprintf(w, "  + %s\n", c.Name)
		case "removed":
			fmt.Fprintf(w, "  - %s\n", c.Name)
		}
	}
}
