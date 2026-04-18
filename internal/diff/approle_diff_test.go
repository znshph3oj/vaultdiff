package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseAppRole = &vault.AppRoleInfo{
	RoleID:       "role-abc",
	BindSecretID: true,
	Policies:     []string{"default", "dev"},
	TTL:          3600,
	MaxTTL:       7200,
}

func TestCompareAppRoles_NoChanges(t *testing.T) {
	b := *baseAppRole
	changes := CompareAppRoles(baseAppRole, &b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareAppRoles_TTLChanged(t *testing.T) {
	b := *baseAppRole
	b.TTL = 1800
	changes := CompareAppRoles(baseAppRole, &b)
	if len(changes) != 1 || changes[0].Field != "token_ttl" {
		t.Errorf("expected token_ttl change, got %+v", changes)
	}
}

func TestCompareAppRoles_PolicyAdded(t *testing.T) {
	b := *baseAppRole
	b.Policies = []string{"default", "dev", "ops"}
	changes := CompareAppRoles(baseAppRole, &b)
	var found bool
	for _, c := range changes {
		if c.Field == "policy_added" && c.NewVal == "ops" {
			found = true
		}
	}
	if !found {
		t.Error("expected policy_added for 'ops'")
	}
}

func TestCompareAppRoles_PolicyRemoved(t *testing.T) {
	b := *baseAppRole
	b.Policies = []string{"default"}
	changes := CompareAppRoles(baseAppRole, &b)
	var found bool
	for _, c := range changes {
		if c.Field == "policy_removed" && c.OldVal == "dev" {
			found = true
		}
	}
	if !found {
		t.Error("expected policy_removed for 'dev'")
	}
}

func TestCompareAppRoles_NilInputs(t *testing.T) {
	if CompareAppRoles(nil, baseAppRole) != nil {
		t.Error("expected nil for nil first arg")
	}
	if CompareAppRoles(baseAppRole, nil) != nil {
		t.Error("expected nil for nil second arg")
	}
}

func TestFprintAppRoleDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintAppRoleDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No AppRole changes") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
