package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func kubernetesServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetKubernetesRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"bound_service_account_names":      []string{"my-sa"},
			"bound_service_account_namespaces": []string{"default"},
			"ttl":      "1h",
			"max_ttl":  "24h",
			"policies": []string{"read-only"},
		},
	}
	srv := kubernetesServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	role, err := GetKubernetesRoleInfo(client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.Name != "my-role" {
		t.Errorf("expected name 'my-role', got %q", role.Name)
	}
	if role.TTL != "1h" {
		t.Errorf("expected TTL '1h', got %q", role.TTL)
	}
	if len(role.Policies) != 1 || role.Policies[0] != "read-only" {
		t.Errorf("unexpected policies: %v", role.Policies)
	}
}

func TestGetKubernetesRoleInfo_NotFound(t *testing.T) {
	srv := kubernetesServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	_, err := GetKubernetesRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetKubernetesRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := kubernetesServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	_, err := GetKubernetesRoleInfo(client, "my-role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
