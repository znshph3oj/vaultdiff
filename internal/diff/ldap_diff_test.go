package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseLDAPRole = &vault.LDAPRoleInfo{
	RoleName:    "dev-group",
	GroupDN:     "ou=groups,dc=example,dc=com",
	UserDN:      "ou=users,dc=example,dc=com",
	BindDN:      "cn=admin,dc=example,dc=com",
	GroupFilter: "(|(memberUid={{.Username}}))",
	GroupAttr:   "cn",
	UserAttr:    "uid",
	TTL:         "1h",
	MaxTTL:      "24h",
	Policies:    []string{"read", "list"},
}

func TestCompareLDAPRoles_NoChanges(t *testing.T) {
	copy := *baseLDAPRole
	if changes := CompareLDAPRoles(baseLDAPRole, &copy); len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareLDAPRoles_TTLChanged(t *testing.T) {
	updated := *baseLDAPRole
	updated.TTL = "2h"
	changes := CompareLDAPRoles(baseLDAPRole, &updated)
	if len(changes) != 1 || changes[0].Field != "ttl" {
		t.Errorf("expected ttl change, got %+v", changes)
	}
}

func TestCompareLDAPRoles_GroupDNChanged(t *testing.T) {
	updated := *baseLDAPRole
	updated.GroupDN = "ou=newgroups,dc=example,dc=com"
	changes := CompareLDAPRoles(baseLDAPRole, &updated)
	if len(changes) != 1 || changes[0].Field != "group_dn" {
		t.Errorf("expected group_dn change, got %+v", changes)
	}
}

func TestCompareLDAPRoles_NilInputs(t *testing.T) {
	if changes := CompareLDAPRoles(nil, baseLDAPRole); changes != nil {
		t.Error("expected nil for nil input")
	}
	if changes := CompareLDAPRoles(baseLDAPRole, nil); changes != nil {
		t.Error("expected nil for nil input")
	}
}

func TestFprintLDAPDiff_NoChanges(t *testing.T) {
	copy := *baseLDAPRole
	var buf bytes.Buffer
	FprintLDAPDiff(&buf, baseLDAPRole, &copy)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-change message, got %q", buf.String())
	}
}

func TestFprintLDAPDiff_ShowsChanges(t *testing.T) {
	updated := *baseLDAPRole
	updated.UserDN = "ou=newusers,dc=example,dc=com"
	var buf bytes.Buffer
	FprintLDAPDiff(&buf, baseLDAPRole, &updated)
	if !strings.Contains(buf.String(), "user_dn") {
		t.Errorf("expected user_dn in output, got %q", buf.String())
	}
}
