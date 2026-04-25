package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseTerraform = &vault.TerraformRoleInfo{
	Organization: "my-org",
	TeamID:       "team-abc",
	TTL:          3600,
	MaxTTL:       7200,
}

func TestCompareTerraformRoles_NoChanges(t *testing.T) {
	other := *baseTerraform
	changes := diff.CompareTerraformRoles(baseTerraform, &other)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareTerraformRoles_OrganizationChanged(t *testing.T) {
	other := *baseTerraform
	other.Organization = "new-org"
	changes := diff.CompareTerraformRoles(baseTerraform, &other)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "organization" {
		t.Errorf("expected field organization, got %s", changes[0].Field)
	}
}

func TestCompareTerraformRoles_TTLChanged(t *testing.T) {
	other := *baseTerraform
	other.TTL = 1800
	changes := diff.CompareTerraformRoles(baseTerraform, &other)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "ttl" {
		t.Errorf("expected field ttl, got %s", changes[0].Field)
	}
}

func TestCompareTerraformRoles_NilInputs(t *testing.T) {
	changes := diff.CompareTerraformRoles(nil, nil)
	if changes != nil {
		t.Errorf("expected nil changes for nil inputs, got %v", changes)
	}
}

func TestFprintTerraformDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	diff.FprintTerraformDiff(&buf, nil)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestFprintTerraformDiff_WithChanges(t *testing.T) {
	other := *baseTerraform
	other.TeamID = "team-xyz"
	changes := diff.CompareTerraformRoles(baseTerraform, &other)
	var buf bytes.Buffer
	diff.FprintTerraformDiff(&buf, changes)
	if !strings.Contains(buf.String(), "team_id") {
		t.Errorf("expected team_id in output, got: %s", buf.String())
	}
}
