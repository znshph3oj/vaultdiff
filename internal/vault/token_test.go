package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func tokenServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestLookupToken_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"accessor":    "abc123",
			"policies":    []string{"default", "admin"},
			"ttl":         3600,
			"expire_time": "2099-01-01T00:00:00Z",
			"renewable":   true,
			"meta":        map[string]string{"user": "alice"},
		},
	}
	srv := tokenServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, _ := vault.NewClient(srv.URL, "test-token", "secret")
	info, err := c.LookupToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Accessor != "abc123" {
		t.Errorf("expected accessor abc123, got %s", info.Accessor)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.TTL != 3600 {
		t.Errorf("expected TTL 3600, got %d", info.TTL)
	}
	if !info.Renewable {
		t.Error("expected renewable to be true")
	}
	if info.ExpireTime.Year() != 2099 {
		t.Errorf("expected expire year 2099, got %d", info.ExpireTime.Year())
	}
}

func TestLookupToken_Unauthorized(t *testing.T) {
	srv := tokenServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c, _ := vault.NewClient(srv.URL, "bad-token", "secret")
	_, err := c.LookupToken()
	if err == nil {
		t.Fatal("expected error for unauthorized token")
	}
}

func TestLookupToken_NoExpireTime(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"accessor":  "xyz",
			"policies":  []string{"default"},
			"ttl":       0,
			"renewable": false,
		},
	}
	srv := tokenServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, _ := vault.NewClient(srv.URL, "test-token", "secret")
	info, err := c.LookupToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.ExpireTime.IsZero() {
		t.Errorf("expected zero expire time, got %v", info.ExpireTime)
	}
	_ = time.Now() // ensure time package is used
}
