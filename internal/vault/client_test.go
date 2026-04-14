package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func mockVaultServer(t *testing.T, version int, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"data": data,
				"metadata": map[string]interface{}{
					"version": float64(version),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestNewClient_DefaultMount(t *testing.T) {
	client, err := vault.NewClient("", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", client.Mount)
	}
}

func TestNewClient_CustomMount(t *testing.T) {
	client, err := vault.NewClient("", "", "kv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", client.Mount)
	}
}

func TestGetSecretVersion_Success(t *testing.T) {
	expectedData := map[string]interface{}{"username": "admin", "password": "s3cr3t"}
	server := mockVaultServer(t, 3, expectedData)
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	sv, err := client.GetSecretVersion(context.Background(), "myapp/config", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sv.Version != 3 {
		t.Errorf("expected version 3, got %d", sv.Version)
	}
	if sv.Data["username"] != "admin" {
		t.Errorf("expected username 'admin', got %v", sv.Data["username"])
	}
}

func TestGetSecretVersion_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetSecretVersion(context.Background(), "missing/secret", 1)
	if err == nil {
		t.Error("expected error for missing secret, got nil")
	}
}
