package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func approleServer(t *testing.T, roleName string, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetAppRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"role_id":         "test-role-id",
			"bind_secret_id":  true,
			"token_policies":  []string{"default", "dev"},
			"token_ttl":       3600,
			"token_max_ttl":   7200,
		},
	}
	srv := approleServer(t, "myrole", http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	info, err := GetAppRoleInfo(client, "myrole")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RoleID != "test-role-id" {
		t.Errorf("expected role_id 'test-role-id', got %q", info.RoleID)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
}

func TestGetAppRoleInfo_NotFound(t *testing.T) {
	srv := approleServer(t, "missing", http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetAppRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAppRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := approleServer(t, "role", http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetAppRoleInfo(client, "role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
