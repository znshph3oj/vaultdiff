package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func samlServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetSAMLRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"bound_subjects":  []string{"user@example.com"},
			"token_policies":  []string{"default", "saml-policy"},
			"token_ttl":       3600,
			"token_max_ttl":   7200,
			"bound_attributes": map[string]string{"department": "engineering"},
		},
	}
	srv := samlServer(t, http.StatusOK, payload)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := GetSAMLRoleInfo(client, "saml", "eng-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "eng-role" {
		t.Errorf("expected name %q, got %q", "eng-role", info.Name)
	}
	if info.TokenTTL != 3600 {
		t.Errorf("expected TTL 3600, got %d", info.TokenTTL)
	}
	if len(info.TokenPolicies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.TokenPolicies))
	}
}

func TestGetSAMLRoleInfo_NotFound(t *testing.T) {
	srv := samlServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetSAMLRoleInfo(client, "saml", "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSAMLRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := samlServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", "")
	_, err := GetSAMLRoleInfo(client, "saml", "role")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
