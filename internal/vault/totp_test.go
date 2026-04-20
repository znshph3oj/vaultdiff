package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/your/vaultdiff/internal/vault"
)

func totpServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetTOTPKeyInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"account_name": "user@example.com",
			"algorithm":    "SHA1",
			"digits":       6,
			"issuer":       "MyApp",
			"period":       30,
			"qr_size":      200,
		},
	}
	server := totpServer(t, http.StatusOK, payload)
	defer server.Close()

	cfg := api.DefaultConfig()
	cfg.Address = server.URL
	client, err := vault.NewClient(cfg, "totp")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := vault.GetTOTPKeyInfo(client, "totp", "mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.AccountName != "user@example.com" {
		t.Errorf("expected account_name %q, got %q", "user@example.com", info.AccountName)
	}
	if info.Period != 30 {
		t.Errorf("expected period 30, got %d", info.Period)
	}
	if info.Digits != 6 {
		t.Errorf("expected digits 6, got %d", info.Digits)
	}
}

func TestGetTOTPKeyInfo_NotFound(t *testing.T) {
	server := totpServer(t, http.StatusNotFound, nil)
	defer server.Close()

	cfg := api.DefaultConfig()
	cfg.Address = server.URL
	client, err := vault.NewClient(cfg, "totp")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = vault.GetTOTPKeyInfo(client, "totp", "missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetTOTPKeyInfo_UnexpectedStatus(t *testing.T) {
	server := totpServer(t, http.StatusInternalServerError, nil)
	defer server.Close()

	cfg := api.DefaultConfig()
	cfg.Address = server.URL
	client, err := vault.NewClient(cfg, "totp")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = vault.GetTOTPKeyInfo(client, "totp", "mykey")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
