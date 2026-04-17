package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultdiff/internal/vault"
)

func makePolicyMap(entries ...vault.PolicyInfo) map[string]*vault.PolicyInfo {
	m := map[string]*vault.PolicyInfo{}
	for i := range entries {
		m[entries[i].Name] = &entries[i]
	}
	return m
}

func TestComparePolicies_Added(t *testing.T) {
	old := makePolicyMap()
	new := makePolicyMap(vault.PolicyInfo{Name: "dev", Rules: "path \"secret/*\" {}"})
	changes := ComparePolicies(old, new)
	if len(changes) != 1 || changes[0].Status != "added" {
		t.Fatalf("expected added, got %+v", changes)
	}
}

func TestComparePolicies_Removed(t *testing.T) {
	old := makePolicyMap(vault.PolicyInfo{Name: "dev", Rules: "path \"secret/*\" {}"})
	new := makePolicyMap()
	changes := ComparePolicies(old, new)
	if len(changes) != 1 || changes[0].Status != "removed" {
		t.Fatalf("expected removed, got %+v", changes)
	}
}

func TestComparePolicies_Modified(t *testing.T) {
	old := makePolicyMap(vault.PolicyInfo{Name: "dev", Rules: "old rules"})
	new := makePolicyMap(vault.PolicyInfo{Name: "dev", Rules: "new rules"})
	changes := ComparePolicies(old, new)
	if len(changes) != 1 || changes[0].Status != "modified" {
		t.Fatalf("expected modified, got %+v", changes)
	}
}

func TestComparePolicies_Unchanged(t *testing.T) {
	old := makePolicyMap(vault.PolicyInfo{Name: "dev", Rules: "same"})
	new := makePolicyMap(vault.PolicyInfo{Name: "dev", Rules: "same"})
	changes := ComparePolicies(old, new)
	if len(changes) != 1 || changes[0].Status != "unchanged" {
		t.Fatalf("expected unchanged, got %+v", changes)
	}
}

func TestFprintPolicyDiff_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	FprintPolicyDiff(&buf, []PolicyChange{{Name: "admin", Status: "added"}})
	if !strings.Contains(buf.String(), "Policy Diff") {
		t.Error("expected header in output")
	}
}
