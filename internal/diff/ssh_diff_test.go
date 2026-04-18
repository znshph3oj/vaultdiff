package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseSSHRole = &vault.SSHRoleInfo{
	KeyType:     "ca",
	DefaultUser: "ubuntu",
	TTL:         "1h",
	MaxTTL:      "24h",
	Policies:    []string{"default"},
}

func TestCompareSSHRoles_NoChanges(t *testing.T) {
	changes := CompareSSHRoles(baseSSHRole, baseSSHRole)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareSSHRoles_KeyTypeChanged(t *testing.T) {
	b := *baseSSHRole
	b.KeyType = "otp"
	changes := CompareSSHRoles(baseSSHRole, &b)
	if len(changes) != 1 || changes[0].Field != "key_type" {
		t.Errorf("expected key_type change, got %+v", changes)
	}
}

func TestCompareSSHRoles_TTLChanged(t *testing.T) {
	b := *baseSSHRole
	b.TTL = "2h"
	changes := CompareSSHRoles(baseSSHRole, &b)
	if len(changes) != 1 || changes[0].Field != "ttl" {
		t.Errorf("expected ttl change, got %+v", changes)
	}
}

func TestCompareSSHRoles_NilInputs(t *testing.T) {
	if CompareSSHRoles(nil, baseSSHRole) != nil {
		t.Error("expected nil for nil a")
	}
	if CompareSSHRoles(baseSSHRole, nil) != nil {
		t.Error("expected nil for nil b")
	}
}

func TestFprintSSHDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintSSHDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No SSH role changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestFprintSSHDiff_WithChanges(t *testing.T) {
	changes := []SSHRoleChange{{Field: "ttl", From: "1h", To: "2h"}}
	var buf bytes.Buffer
	FprintSSHDiff(&buf, changes)
	if !strings.Contains(buf.String(), "ttl") {
		t.Errorf("expected ttl in output, got: %s", buf.String())
	}
}
