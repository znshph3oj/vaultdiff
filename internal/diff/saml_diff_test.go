package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseSAML = &vault.SAMLRoleInfo{
	Name:          "eng-role",
	BoundSubjects: []string{"user@example.com"},
	TokenPolicies: []string{"default"},
	TokenTTL:      3600,
	TokenMaxTTL:   7200,
}

func TestCompareSAMLRoles_NoChanges(t *testing.T) {
	other := *baseSAML
	changes := CompareSAMLRoles(baseSAML, &other)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareSAMLRoles_TTLChanged(t *testing.T) {
	other := *baseSAML
	other.TokenTTL = 1800
	changes := CompareSAMLRoles(baseSAML, &other)
	if len(changes) != 1 || changes[0].Field != "token_ttl" {
		t.Errorf("expected token_ttl change, got %+v", changes)
	}
}

func TestCompareSAMLRoles_PolicyChanged(t *testing.T) {
	other := *baseSAML
	other.TokenPolicies = []string{"default", "saml-extra"}
	changes := CompareSAMLRoles(baseSAML, &other)
	if len(changes) != 1 || changes[0].Field != "token_policies" {
		t.Errorf("expected token_policies change, got %+v", changes)
	}
}

func TestCompareSAMLRoles_SubjectsChanged(t *testing.T) {
	other := *baseSAML
	other.BoundSubjects = []string{"admin@example.com"}
	changes := CompareSAMLRoles(baseSAML, &other)
	if len(changes) != 1 || changes[0].Field != "bound_subjects" {
		t.Errorf("expected bound_subjects change, got %+v", changes)
	}
}

func TestCompareSAMLRoles_NilInputs(t *testing.T) {
	changes := CompareSAMLRoles(nil, nil)
	if changes != nil {
		t.Errorf("expected nil changes for nil inputs")
	}
}

func TestFprintSAMLDiff_NoChanges(t *testing.T) {
	other := *baseSAML
	var buf bytes.Buffer
	FprintSAMLDiff(&buf, baseSAML, &other)
	if !strings.Contains(buf.String(), "No SAML role changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestFprintSAMLDiff_WithChanges(t *testing.T) {
	other := *baseSAML
	other.TokenTTL = 900
	var buf bytes.Buffer
	FprintSAMLDiff(&buf, baseSAML, &other)
	if !strings.Contains(buf.String(), "token_ttl") {
		t.Errorf("expected token_ttl in output, got: %s", buf.String())
	}
}
