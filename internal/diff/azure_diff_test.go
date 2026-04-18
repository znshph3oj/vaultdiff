package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseAzureRole = &vault.AzureRoleInfo{
	ApplicationObjectID: "app-123",
	ClientID:            "client-456",
	TTL:                 "1h",
	MaxTTL:              "24h",
	AzureRoles:          []string{"Contributor"},
	AzureGroups:         []string{"devs"},
}

func TestCompareAzureRoles_NoChanges(t *testing.T) {
	changes := CompareAzureRoles(baseAzureRole, baseAzureRole)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareAzureRoles_TTLChanged(t *testing.T) {
	b := *baseAzureRole
	b.TTL = "2h"
	changes := CompareAzureRoles(baseAzureRole, &b)
	if len(changes) != 1 || changes[0].Field != "ttl" {
		t.Errorf("expected ttl change, got %+v", changes)
	}
}

func TestCompareAzureRoles_RolesChanged(t *testing.T) {
	b := *baseAzureRole
	b.AzureRoles = []string{"Reader"}
	changes := CompareAzureRoles(baseAzureRole, &b)
	if len(changes) != 1 || changes[0].Field != "azure_roles" {
		t.Errorf("expected azure_roles change, got %+v", changes)
	}
}

func TestCompareAzureRoles_NilInputs(t *testing.T) {
	if CompareAzureRoles(nil, baseAzureRole) != nil {
		t.Error("expected nil for nil input a")
	}
	if CompareAzureRoles(baseAzureRole, nil) != nil {
		t.Error("expected nil for nil input b")
	}
}

func TestFprintAzureDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintAzureDiff(&buf, baseAzureRole, baseAzureRole)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestFprintAzureDiff_ShowsChanges(t *testing.T) {
	b := *baseAzureRole
	b.ClientID = "new-client"
	var buf bytes.Buffer
	FprintAzureDiff(&buf, baseAzureRole, &b)
	if !strings.Contains(buf.String(), "client_id") {
		t.Errorf("expected client_id in output, got: %s", buf.String())
	}
}
