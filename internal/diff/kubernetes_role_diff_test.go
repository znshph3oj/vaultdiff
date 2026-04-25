package diff

import (
	"bytes"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

var baseKubernetesRole = &vault.KubernetesRoleBinding{
	Name:                          "my-role",
	BoundServiceAccountNames:      []string{"default"},
	BoundServiceAccountNamespaces: []string{"default"},
	TTL:                           "1h",
	MaxTTL:                        "24h",
	Policies:                      []string{"read-only"},
}

func TestCompareKubernetesRoles_NoChanges(t *testing.T) {
	changes := CompareKubernetesRoles(baseKubernetesRole, baseKubernetesRole)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareKubernetesRoles_TTLChanged(t *testing.T) {
	b := *baseKubernetesRole
	b.TTL = "2h"
	changes := CompareKubernetesRoles(baseKubernetesRole, &b)
	if len(changes) != 1 || changes[0].Field != "ttl" {
		t.Errorf("expected ttl change, got %+v", changes)
	}
}

func TestCompareKubernetesRoles_PolicyAdded(t *testing.T) {
	b := *baseKubernetesRole
	b.Policies = []string{"read-only", "write-secrets"}
	changes := CompareKubernetesRoles(baseKubernetesRole, &b)
	if len(changes) != 1 || changes[0].Field != "token_policies" {
		t.Errorf("expected token_policies change, got %+v", changes)
	}
}

func TestCompareKubernetesRoles_NamespaceChanged(t *testing.T) {
	b := *baseKubernetesRole
	b.BoundServiceAccountNamespaces = []string{"kube-system"}
	changes := CompareKubernetesRoles(baseKubernetesRole, &b)
	if len(changes) != 1 || changes[0].Field != "bound_service_account_namespaces" {
		t.Errorf("expected namespace change, got %+v", changes)
	}
}

func TestCompareKubernetesRoles_NilInputs(t *testing.T) {
	if changes := CompareKubernetesRoles(nil, nil); len(changes) != 0 {
		t.Errorf("expected no changes for nil inputs, got %d", len(changes))
	}
}

func TestFprintKubernetesRoleDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintKubernetesRoleDiff(&buf, baseKubernetesRole, baseKubernetesRole)
	if buf.String() != "No changes in Kubernetes role binding.\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFprintKubernetesRoleDiff_WithChanges(t *testing.T) {
	b := *baseKubernetesRole
	b.TTL = "4h"
	var buf bytes.Buffer
	FprintKubernetesRoleDiff(&buf, baseKubernetesRole, &b)
	if buf.Len() == 0 {
		t.Error("expected non-empty diff output")
	}
}
