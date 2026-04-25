package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func baseOIDC() *vault.OIDCRoleInfo {
	return &vault.OIDCRoleInfo{
		RoleName:         "myrole",
		BoundAudiences:   []string{"vault"},
		AllowedRedirects: []string{"https://example.com/callback"},
		UserClaim:        "sub",
		TokenTTL:         3600,
		TokenMaxTTL:      7200,
		TokenPolicies:    []string{"default"},
	}
}

func TestCompareOIDCRoles_NoChanges(t *testing.T) {
	a := baseOIDC()
	b := baseOIDC()
	changes := CompareOIDCRoles(a, b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareOIDCRoles_TTLChanged(t *testing.T) {
	a := baseOIDC()
	b := baseOIDC()
	b.TokenTTL = 1800
	changes := CompareOIDCRoles(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "token_ttl" {
		t.Errorf("expected field token_ttl, got %s", changes[0].Field)
	}
}

func TestCompareOIDCRoles_PolicyChanged(t *testing.T) {
	a := baseOIDC()
	b := baseOIDC()
	b.TokenPolicies = []string{"default", "admin"}
	changes := CompareOIDCRoles(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "token_policies" {
		t.Errorf("expected field token_policies, got %s", changes[0].Field)
	}
}

func TestCompareOIDCRoles_NilInputs(t *testing.T) {
	changes := CompareOIDCRoles(nil, nil)
	if changes != nil {
		t.Errorf("expected nil for nil inputs, got %v", changes)
	}
}

func TestFprintOIDCDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintOIDCDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No OIDC role changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestFprintOIDCDiff_WithChanges(t *testing.T) {
	changes := []OIDCRoleChange{
		{Field: "user_claim", From: "sub", To: "email"},
	}
	var buf bytes.Buffer
	FprintOIDCDiff(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "user_claim") {
		t.Errorf("expected user_claim in output, got: %s", out)
	}
	if !strings.Contains(out, "email") {
		t.Errorf("expected 'email' in output, got: %s", out)
	}
}
