package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ldapServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetLDAPRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"groupdn":  "ou=groups,dc=example,dc=com",
			"userdn":   "ou=users,dc=example,dc=com",
			"ttl":      "1h",
			"policies": []string{"read"},
		},
	}
	srv := ldapServer(t, http.StatusOK, payload)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatal(err)
	}

	info, err := GetLDAPRoleInfo(client, "ldap", "dev-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.GroupDN != "ou=groups,dc=example,dc=com" {
		t.Errorf("expected groupdn, got %q", info.GroupDN)
	}
	if info.RoleName != "dev-group" {
		t.Errorf("expected role name dev-group, got %q", info.RoleName)
	}
}

func TestGetLDAPRoleInfo_NotFound(t *testing.T) {
	srv := ldapServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetLDAPRoleInfo(client, "ldap", "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLDAPRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := ldapServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetLDAPRoleInfo(client, "ldap", "dev-group")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLDAPRoleInfo_EmptyRoleName(t *testing.T) {
	srv := ldapServer(t, http.StatusOK, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetLDAPRoleInfo(client, "ldap", "")
	if err == nil {
		t.Fatal("expected error for empty role name, got nil")
	}
}
