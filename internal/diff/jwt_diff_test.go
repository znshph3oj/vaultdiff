package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func baseJWT() *vault.JWTRoleInfo {
	return &vault.JWTRoleInfo{
		Name:           "test-role",
		RoleType:       "jwt",
		BoundAudiences: []string{"https://example.com"},
		UserClaim:      "sub",
		GroupsClaim:    "groups",
		TTL:            "1h",
		MaxTTL:         "24h",
		TokenPolicies:  []string{"default"},
	}
}

func TestCompareJWTRoles_NoChanges(t *testing.T) {
	a := baseJWT()
	b := baseJWT()
	changes := CompareJWTRoles(a, b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareJWTRoles_TTLChanged(t *testing.T) {
	a := baseJWT()
	b := baseJWT()
	b.TTL = "2h"
	changes := CompareJWTRoles(a, b)
	if len(changes) != 1 || changes[0].Field != "ttl" {
		t.Errorf("expected ttl change, got %+v", changes)
	}
}

func TestCompareJWTRoles_PolicyChanged(t *testing.T) {
	a := baseJWT()
	b := baseJWT()
	b.TokenPolicies = []string{"default", "admin"}
	changes := CompareJWTRoles(a, b)
	if len(changes) != 1 || changes[0].Field != "token_policies" {
		t.Errorf("expected token_policies change, got %+v", changes)
	}
}

func TestCompareJWTRoles_NilInputs(t *testing.T) {
	if changes := CompareJWTRoles(nil, baseJWT()); changes != nil {
		t.Errorf("expected nil for nil a")
	}
	if changes := CompareJWTRoles(baseJWT(), nil); changes != nil {
		t.Errorf("expected nil for nil b")
	}
}

func TestFprintJWTDiff_NoChanges(t *testing.T) {
	a := baseJWT()
	var buf bytes.Buffer
	FprintJWTDiff(&buf, a, a)
	if !strings.Contains(buf.String(), "No JWT role changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestFprintJWTDiff_ShowsChanges(t *testing.T) {
	a := baseJWT()
	b := baseJWT()
	b.RoleType = "oidc"
	var buf bytes.Buffer
	FprintJWTDiff(&buf, a, b)
	if !strings.Contains(buf.String(), "role_type") {
		t.Errorf("expected role_type in diff output, got: %s", buf.String())
	}
}
