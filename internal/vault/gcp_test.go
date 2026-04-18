package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func gcpServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetGCPRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"name":        "my-role",
			"role_type":   "service_account_key",
			"project":     "my-project",
			"secret_type": "service_account_key",
			"token_ttl":   3600,
		},
	}
	srv := gcpServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	info, err := GetGCPRoleInfo(client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "my-role" {
		t.Errorf("expected name my-role, got %s", info.Name)
	}
	if info.Project != "my-project" {
		t.Errorf("expected project my-project, got %s", info.Project)
	}
}

func TestGetGCPRoleInfo_NotFound(t *testing.T) {
	srv := gcpServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetGCPRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetGCPRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := gcpServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetGCPRoleInfo(client, "my-role")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
