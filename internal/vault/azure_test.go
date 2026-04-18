package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func azureServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetAzureRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"application_object_id": "app-123",
			"client_id":             "client-456",
			"ttl":                   "1h",
			"max_ttl":               "24h",
			"azure_roles":           []string{"Contributor"},
			"azure_groups":          []string{"devs"},
		},
	}
	srv := azureServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	info, err := GetAzureRoleInfo(client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ApplicationObjectID != "app-123" {
		t.Errorf("expected app-123, got %s", info.ApplicationObjectID)
	}
	if len(info.AzureRoles) != 1 || info.AzureRoles[0] != "Contributor" {
		t.Errorf("unexpected azure roles: %v", info.AzureRoles)
	}
}

func TestGetAzureRoleInfo_NotFound(t *testing.T) {
	srv := azureServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetAzureRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAzureRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := azureServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetAzureRoleInfo(client, "my-role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
