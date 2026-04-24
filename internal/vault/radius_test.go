package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func radiusServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetRADIUSRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"policies": []string{"default", "radius-user"},
			"ttl":      "1h",
			"max_ttl":  "24h",
		},
	}
	srv := radiusServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	info, err := GetRADIUSRoleInfo(client, "ops")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.TTL != "1h" {
		t.Errorf("expected TTL '1h', got %q", info.TTL)
	}
	if info.MaxTTL != "24h" {
		t.Errorf("expected MaxTTL '24h', got %q", info.MaxTTL)
	}
}

func TestGetRADIUSRoleInfo_NotFound(t *testing.T) {
	srv := radiusServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetRADIUSRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetRADIUSRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := radiusServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetRADIUSRoleInfo(client, "ops")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
