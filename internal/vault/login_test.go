package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func loginServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestLogin_Success(t *testing.T) {
	payload := map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   "s.abc123",
			"accessor":       "acc-xyz",
			"policies":       []string{"default", "dev"},
			"lease_duration": 3600,
			"renewable":      true,
			"metadata":       map[string]string{"username": "alice"},
		},
	}
	srv := loginServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, err := NewClient(srv.URL, "root", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := c.Login("userpass", map[string]interface{}{"password": "hunter2"})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if info.ClientToken != "s.abc123" {
		t.Errorf("ClientToken = %q, want s.abc123", info.ClientToken)
	}
	if !info.Renewable {
		t.Error("expected Renewable=true")
	}
	if info.LeaseDuration.Seconds() != 3600 {
		t.Errorf("LeaseDuration = %v, want 3600s", info.LeaseDuration)
	}
	if info.Meta["username"] != "alice" {
		t.Errorf("Meta username = %q, want alice", info.Meta["username"])
	}
}

func TestLogin_NotFound(t *testing.T) {
	srv := loginServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "root", "secret")
	_, err := c.Login("unknown", nil)
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestLogin_Forbidden(t *testing.T) {
	srv := loginServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "root", "secret")
	_, err := c.Login("userpass", map[string]interface{}{"password": "wrong"})
	if err == nil {
		t.Fatal("expected error for 403")
	}
}
