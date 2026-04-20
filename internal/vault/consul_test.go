package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func consulServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetConsulRoleInfo_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"policies":   []string{"read-only", "kv-access"},
			"token_type": "client",
			"ttl":        "1h",
			"max_ttl":    "24h",
			"local":      false,
		},
	}
	srv := consulServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := vault.NewClient(srv.URL, "test-token", "secret")
	info, err := vault.GetConsulRoleInfo(client, "consul", "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "my-role" {
		t.Errorf("expected name %q, got %q", "my-role", info.Name)
	}
	if info.TokenType != "client" {
		t.Errorf("expected token_type %q, got %q", "client", info.TokenType)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
}

func TestGetConsulRoleInfo_NotFound(t *testing.T) {
	srv := consulServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := vault.NewClient(srv.URL, "test-token", "secret")
	_, err := vault.GetConsulRoleInfo(client, "consul", "missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetConsulRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := consulServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := vault.NewClient(srv.URL, "test-token", "secret")
	_, err := vault.GetConsulRoleInfo(client, "consul", "my-role")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
