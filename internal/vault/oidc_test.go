package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func oidcServer(t *testing.T, role string, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/auth/oidc/role/" + role
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetOIDCRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"bound_audiences":       []string{"vault"},
			"allowed_redirect_uris": []string{"https://example.com/oidc/callback"},
			"user_claim":            "sub",
			"token_ttl":             3600,
			"token_max_ttl":         7200,
			"token_policies":        []string{"default", "dev"},
		},
	}
	srv := oidcServer(t, "myrole", http.StatusOK, payload)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	info, err := GetOIDCRoleInfo(client, "oidc", "myrole")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.UserClaim != "sub" {
		t.Errorf("expected user_claim=sub, got %s", info.UserClaim)
	}
	if info.TokenTTL != 3600 {
		t.Errorf("expected token_ttl=3600, got %d", info.TokenTTL)
	}
	if len(info.TokenPolicies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.TokenPolicies))
	}
	if info.RoleName != "myrole" {
		t.Errorf("expected role_name=myrole, got %s", info.RoleName)
	}
}

func TestGetOIDCRoleInfo_NotFound(t *testing.T) {
	srv := oidcServer(t, "myrole", http.StatusNotFound, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetOIDCRoleInfo(client, "oidc", "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetOIDCRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := oidcServer(t, "myrole", http.StatusInternalServerError, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetOIDCRoleInfo(client, "oidc", "myrole")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
