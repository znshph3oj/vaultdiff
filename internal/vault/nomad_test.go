package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func nomadServer(t *testing.T, role string, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/nomad/role/" + role
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetNomadRoleInfo_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"policies": []string{"dev", "ops"},
			"global":   false,
			"type":     "client",
			"lease":    "1h",
		},
	}
	srv := nomadServer(t, "myrole", http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	info, err := GetNomadRoleInfo(client, "myrole")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "myrole" {
		t.Errorf("expected name %q, got %q", "myrole", info.Name)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.Type != "client" {
		t.Errorf("expected type %q, got %q", "client", info.Type)
	}
	if info.Lease != "1h" {
		t.Errorf("expected lease %q, got %q", "1h", info.Lease)
	}
}

func TestGetNomadRoleInfo_NotFound(t *testing.T) {
	srv := nomadServer(t, "myrole", http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	_, err := GetNomadRoleInfo(client, "myrole")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestGetNomadRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := nomadServer(t, "myrole", http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	_, err := GetNomadRoleInfo(client, "myrole")
	if err == nil {
		t.Fatal("expected error for unexpected status")
	}
}
