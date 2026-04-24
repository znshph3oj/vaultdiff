package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestCompareGitHubRoles_NoChanges(t *testing.T) {
	a := &vault.GitHubRoleInfo{TeamName: "devops", Policies: []string{"read"}, TTL: "1h", MaxTTL: "24h"}
	b := &vault.GitHubRoleInfo{TeamName: "devops", Policies: []string{"read"}, TTL: "1h", MaxTTL: "24h"}
	changes := CompareGitHubRoles(a, b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareGitHubRoles_TTLChanged(t *testing.T) {
	a := &vault.GitHubRoleInfo{TTL: "1h", MaxTTL: "24h"}
	b := &vault.GitHubRoleInfo{TTL: "2h", MaxTTL: "24h"}
	changes := CompareGitHubRoles(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "ttl" || changes[0].From != "1h" || changes[0].To != "2h" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompareGitHubRoles_PolicyAdded(t *testing.T) {
	a := &vault.GitHubRoleInfo{Policies: []string{"read"}}
	b := &vault.GitHubRoleInfo{Policies: []string{"read", "write"}}
	changes := CompareGitHubRoles(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "policy" || changes[0].To != "write" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompareGitHubRoles_PolicyRemoved(t *testing.T) {
	a := &vault.GitHubRoleInfo{Policies: []string{"read", "admin"}}
	b := &vault.GitHubRoleInfo{Policies: []string{"read"}}
	changes := CompareGitHubRoles(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "policy" || changes[0].From != "admin" || changes[0].To != "" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompareGitHubRoles_NilInputs(t *testing.T) {
	changes := CompareGitHubRoles(nil, nil)
	if changes != nil {
		t.Errorf("expected nil for nil inputs")
	}
}

func TestFprintGitHubDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintGitHubDiff(&buf, "devops", nil)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestFprintGitHubDiff_WithChanges(t *testing.T) {
	changes := []GitHubRoleChange{
		{Field: "ttl", From: "1h", To: "2h"},
		{Field: "policy", From: "", To: "write"},
	}
	var buf bytes.Buffer
	FprintGitHubDiff(&buf, "devops", changes)
	out := buf.String()
	if !strings.Contains(out, "~ ttl") {
		t.Errorf("expected ttl change line, got: %s", out)
	}
	if !strings.Contains(out, "+ policy") {
		t.Errorf("expected policy addition line, got: %s", out)
	}
}
