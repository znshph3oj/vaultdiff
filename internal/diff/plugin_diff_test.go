package diff

import (
	"bytes"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestComparePlugins_NoChanges(t *testing.T) {
	old := &vault.PluginList{Plugins: []string{"plugin-a", "plugin-b"}}
	new := &vault.PluginList{Plugins: []string{"plugin-a", "plugin-b"}}

	result := ComparePlugins(old, new)
	if len(result.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestComparePlugins_Added(t *testing.T) {
	old := &vault.PluginList{Plugins: []string{"plugin-a"}}
	new := &vault.PluginList{Plugins: []string{"plugin-a", "plugin-b"}}

	result := ComparePlugins(old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Name != "plugin-b" || result.Changes[0].Status != "added" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestComparePlugins_Removed(t *testing.T) {
	old := &vault.PluginList{Plugins: []string{"plugin-a", "plugin-b"}}
	new := &vault.PluginList{Plugins: []string{"plugin-a"}}

	result := ComparePlugins(old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Name != "plugin-b" || result.Changes[0].Status != "removed" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestComparePlugins_NilInputs(t *testing.T) {
	result := ComparePlugins(nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected no changes for nil inputs, got %d", len(result.Changes))
	}
}

func TestFprintPluginDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintPluginDiff(&buf, &PluginDiffResult{})
	if buf.String() == "" {
		t.Error("expected non-empty output")
	}
	if buf.String() != "No plugin changes detected.\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFprintPluginDiff_WithChanges(t *testing.T) {
	result := &PluginDiffResult{
		Changes: []PluginChange{
			{Name: "plugin-a", Status: "added"},
			{Name: "plugin-b", Status: "removed"},
		},
	}
	var buf bytes.Buffer
	FprintPluginDiff(&buf, result)
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("+ plugin-a")) {
		t.Errorf("expected added plugin in output, got: %s", out)
	}
	if !bytes.Contains(buf.Bytes(), []byte("- plugin-b")) {
		t.Errorf("expected removed plugin in output, got: %s", out)
	}
}
