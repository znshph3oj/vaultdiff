package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func jwtServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetJWTRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"role_type":       "jwt",
			"bound_audiences": []string{"https://example.com"},
			"user_claim":      "sub",
			"groups_claim":    "groups",
			"ttl":             "1h",
			"max_ttl":         "24h",
			"token_policies":  []string{"default", "dev"},
		},
	}
	srv := jwtServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := &Client{addr: srv.URL, http: srv.Client()}
	info, err := c.GetJWTRoleInfo("jwt", "dev-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RoleType != "jwt" {
		t.Errorf("expected role_type jwt, got %s", info.RoleType)
	}
	if info.Name != "dev-role" {
		t.Errorf("expected name dev-role, got %s", info.Name)
	}
	if len(info.TokenPolicies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.TokenPolicies))
	}
}

func TestGetJWTRoleInfo_NotFound(t *testing.T) {
	srv := jwtServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := &Client{addr: srv.URL, http: srv.Client()}
	_, err := c.GetJWTRoleInfo("jwt", "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetJWTRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := jwtServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := &Client{addr: srv.URL, http: srv.Client()}
	_, err := c.GetJWTRoleInfo("jwt", "some-role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
