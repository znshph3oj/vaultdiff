package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func baseMounts() map[string]*vault.MountInfo {
	return map[string]*vault.MountInfo{
		"secret/": {Path: "secret/", Type: "kv", Description: "kv store"},
		"pki/":    {Path: "pki/", Type: "pki", Description: "PKI"},
	}
}

func TestCompareMounts_Added(t *testing.T) {
	before := baseMounts()
	after := baseMounts()
	after["transit/"] = &vault.MountInfo{Path: "transit/", Type: "transit"}

	changes := CompareMounts(before, after)
	if len(changes) != 1 || changes[0].Status != "added" || changes[0].Path != "transit/" {
		t.Errorf("expected one added change for transit/, got %+v", changes)
	}
}

func TestCompareMounts_Removed(t *testing.T) {
	before := baseMounts()
	after := baseMounts()
	delete(after, "pki/")

	changes := CompareMounts(before, after)
	if len(changes) != 1 || changes[0].Status != "removed" || changes[0].Path != "pki/" {
		t.Errorf("expected one removed change for pki/, got %+v", changes)
	}
}

func TestCompareMounts_Modified(t *testing.T) {
	before := baseMounts()
	after := baseMounts()
	after["secret/"].Type = "kv-v2"

	changes := CompareMounts(before, after)
	if len(changes) != 1 || changes[0].Status != "modified" {
		t.Errorf("expected one modified change, got %+v", changes)
	}
}

func TestCompareMounts_NoChanges(t *testing.T) {
	before := baseMounts()
	after := baseMounts()

	changes := CompareMounts(before, after)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestFprintMountDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintMountDiff(&buf, []MountChange{})
	if !strings.Contains(buf.String(), "No mount changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestFprintMountDiff_ShowsChanges(t *testing.T) {
	changes := []MountChange{
		{Path: "transit/", Status: "added", New: &vault.MountInfo{Type: "transit"}},
		{Path: "pki/", Status: "removed", Old: &vault.MountInfo{Type: "pki"}},
	}
	var buf bytes.Buffer
	FprintMountDiff(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "+ transit/") {
		t.Errorf("expected added transit/ in output, got: %s", out)
	}
	if !strings.Contains(out, "- pki/") {
		t.Errorf("expected removed pki/ in output, got: %s", out)
	}
}
