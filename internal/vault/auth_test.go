package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func authServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetAuthInfo_Success(t *testing.T) {
	body := map[string]interface{}{
		"data": map[string]interface{}{
			"accessor":     "abc123",
			"display_name": "token-user",
			"policies":     []string{"default", "admin"},
			"ttl":          3600,
			"renewable":    true,
			"meta":         map[string]string{"env": "prod"},
		},
	}
	srv := authServer(t, http.StatusOK, body)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	info, err := c.GetAuthInfo()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Accessor != "abc123" {
		t.Errorf("expected accessor abc123, got %s", info.Accessor)
	}
	if info.DisplayName != "token-user" {
		t.Errorf("expected display_name token-user, got %s", info.DisplayName)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if !info.Renewable {
		t.Error("expected renewable to be true")
	}
	if info.Meta["env"] != "prod" {
		t.Errorf("expected meta env=prod, got %s", info.Meta["env"])
	}
}

func TestGetAuthInfo_Unauthorized(t *testing.T) {
	srv := authServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "bad-token", HTTP: srv.Client()}
	_, err := c.GetAuthInfo()
	if err == nil {
		t.Fatal("expected error for unauthorized")
	}
}

func TestGetAuthInfo_UnexpectedStatus(t *testing.T) {
	srv := authServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "token", HTTP: srv.Client()}
	_, err := c.GetAuthInfo()
	if err == nil {
		t.Fatal("expected error for unexpected status")
	}
}
