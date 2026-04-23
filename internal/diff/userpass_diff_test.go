package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func baseUserpass() *vault.UserpassRoleInfo {
	return &vault.UserpassRoleInfo{
		Username: "alice",
		Policies: []string{"default", "dev"},
		TTL:      "1h",
		MaxTTL:   "24h",
		BoundCIDRs: []string{"10.0.0.0/8"},
	}
}

func TestCompareUserpassRoles_NoChanges(t *testing.T) {
	a := baseUserpass()
	b := baseUserpass()
	changes := CompareUserpassRoles(a, b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareUserpassRoles_TTLChanged(t *testing.T) {
	a := baseUserpass()
	b := baseUserpass()
	b.TTL = "2h"
	changes := CompareUserpassRoles(a, b)
	if len(changes) != 1 || changes[0].Field != "token_ttl" {
		t.Errorf("expected token_ttl change, got %+v", changes)
	}
}

func TestCompareUserpassRoles_PolicyAdded(t *testing.T) {
	a := baseUserpass()
	b := baseUserpass()
	b.Policies = append(b.Policies, "ops")
	changes := CompareUserpassRoles(a, b)
	if len(changes) != 1 || changes[0].Field != "policy_added" || changes[0].NewValue != "ops" {
		t.Errorf("expected policy_added change, got %+v", changes)
	}
}

func TestCompareUserpassRoles_PolicyRemoved(t *testing.T) {
	a := baseUserpass()
	b := baseUserpass()
	b.Policies = []string{"default"}
	changes := CompareUserpassRoles(a, b)
	if len(changes) != 1 || changes[0].Field != "policy_removed" {
		t.Errorf("expected policy_removed change, got %+v", changes)
	}
}

func TestCompareUserpassRoles_NilInputs(t *testing.T) {
	if changes := CompareUserpassRoles(nil, baseUserpass()); changes != nil {
		t.Error("expected nil for nil input a")
	}
	if changes := CompareUserpassRoles(baseUserpass(), nil); changes != nil {
		t.Error("expected nil for nil input b")
	}
}

func TestFprintUserpassDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintUserpassDiff(&buf, "alice", nil)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}

func TestFprintUserpassDiff_WithChanges(t *testing.T) {
	changes := []UserpassChange{
		{Field: "token_ttl", OldValue: "1h", NewValue: "2h"},
		{Field: "policy_added", OldValue: "", NewValue: "ops"},
	}
	var buf bytes.Buffer
	FprintUserpassDiff(&buf, "alice", changes)
	out := buf.String()
	if !strings.Contains(out, "token_ttl") {
		t.Errorf("expected token_ttl in output, got: %s", out)
	}
	if !strings.Contains(out, "ops") {
		t.Errorf("expected ops policy in output, got: %s", out)
	}
}
