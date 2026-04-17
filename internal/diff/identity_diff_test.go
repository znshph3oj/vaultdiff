package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func baseEntity() *vault.IdentityEntity {
	return &vault.IdentityEntity{
		ID:       "abc-123",
		Name:     "alice",
		Policies: []string{"default", "admin"},
		Disabled: false,
	}
}

func TestCompareIdentity_NoChanges(t *testing.T) {
	a, b := baseEntity(), baseEntity()
	if changes := CompareIdentity(a, b); len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareIdentity_NameChanged(t *testing.T) {
	a, b := baseEntity(), baseEntity()
	b.Name = "bob"
	changes := CompareIdentity(a, b)
	if len(changes) != 1 || changes[0].Field != "name" {
		t.Errorf("expected name change, got %+v", changes)
	}
}

func TestCompareIdentity_PolicyAdded(t *testing.T) {
	a, b := baseEntity(), baseEntity()
	b.Policies = append(b.Policies, "superuser")
	changes := CompareIdentity(a, b)
	if len(changes) != 1 || changes[0].NewValue != "superuser" {
		t.Errorf("expected policy addition, got %+v", changes)
	}
}

func TestCompareIdentity_PolicyRemoved(t *testing.T) {
	a, b := baseEntity(), baseEntity()
	b.Policies = []string{"default"}
	changes := CompareIdentity(a, b)
	if len(changes) != 1 || changes[0].OldValue != "admin" {
		t.Errorf("expected policy removal, got %+v", changes)
	}
}

func TestCompareIdentity_NilInputs(t *testing.T) {
	if changes := CompareIdentity(nil, baseEntity()); changes != nil {
		t.Errorf("expected nil for nil input")
	}
}

func TestFprintIdentityDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintIdentityDiff(&buf, baseEntity(), baseEntity())
	if !strings.Contains(buf.String(), "No identity changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestFprintIdentityDiff_ShowsChanges(t *testing.T) {
	a, b := baseEntity(), baseEntity()
	b.Disabled = true
	var buf bytes.Buffer
	FprintIdentityDiff(&buf, a, b)
	if !strings.Contains(buf.String(), "disabled") {
		t.Errorf("expected disabled field in output, got: %s", buf.String())
	}
}
