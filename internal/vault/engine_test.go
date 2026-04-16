package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func engineServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListEngines_Success(t *testing.T) {
	payload := map[string]interface{}{
		"secret/": map[string]interface{}{
			"type":        "kv",
			"description": "key/value secrets",
			"options":     map[string]string{"version": "2"},
			"local":       false,
			"seal_wrap":   false,
		},
		"pki/": map[string]interface{}{
			"type":        "pki",
			"description": "PKI engine",
			"options":     map[string]string{},
			"local":       true,
			"seal_wrap":   true,
		},
	}

	srv := engineServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	engines, err := c.ListEngines(context.Background())
	if err != nil {
		t.Fatalf("ListEngines: %v", err)
	}
	if len(engines) != 2 {
		t.Fatalf("expected 2 engines, got %d", len(engines))
	}
}

func TestListEngines_NotFound(t *testing.T) {
	srv := engineServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.ListEngines(context.Background())
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestListEngines_UnexpectedStatus(t *testing.T) {
	srv := engineServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.ListEngines(context.Background())
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
