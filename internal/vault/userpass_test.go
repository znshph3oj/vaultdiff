package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func userpassServer(t *testing.T, username string, status int, data interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/auth/userpass/users/" + username
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if data != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
		}
	}))
}

func TestGetUserpassRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"token_policies": []string{"default", "dev"},
		"token_ttl":      "1h",
		"token_max_ttl":  "24h",
	}
	srv := userpassServer(t, "alice", http.StatusOK, payload)
	defer srv.Close()

	client := testClient(t, srv.URL)
	info, err := GetUserpassRoleInfo(client, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Username != "alice" {
		t.Errorf("expected username alice, got %s", info.Username)
	}
	if info.TTL != "1h" {
		t.Errorf("expected TTL 1h, got %s", info.TTL)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
}

func TestGetUserpassRoleInfo_NotFound(t *testing.T) {
	srv := userpassServer(t, "alice", http.StatusNotFound, nil)
	defer srv.Close()

	client := testClient(t, srv.URL)
	_, err := GetUserpassRoleInfo(client, "alice")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserpassRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := userpassServer(t, "alice", http.StatusInternalServerError, nil)
	defer srv.Close()

	client := testClient(t, srv.URL)
	_, err := GetUserpassRoleInfo(client, "alice")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
