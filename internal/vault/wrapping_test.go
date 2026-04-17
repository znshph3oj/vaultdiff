package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func wrappingServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/wrapping/lookup" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestLookupWrappingToken_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"token":           "wraps.abc123",
			"accessor":        "acc-xyz",
			"ttl":             300,
			"creation_time":   "2024-01-01T00:00:00Z",
			"creation_path":   "secret/data/myapp",
			"wrapped_accessor": "wacc-xyz",
		},
	}
	srv := wrappingServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "root", HTTP: srv.Client()}
	info, err := c.LookupWrappingToken("wraps.abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Token != "wraps.abc123" {
		t.Errorf("expected token wraps.abc123, got %s", info.Token)
	}
	if info.TTL != 300 {
		t.Errorf("expected TTL 300, got %d", info.TTL)
	}
	if info.CreationPath != "secret/data/myapp" {
		t.Errorf("unexpected creation path: %s", info.CreationPath)
	}
}

func TestLookupWrappingToken_NotFound(t *testing.T) {
	srv := wrappingServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "root", HTTP: srv.Client()}
	_, err := c.LookupWrappingToken("wraps.expired")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLookupWrappingToken_UnexpectedStatus(t *testing.T) {
	srv := wrappingServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "root", HTTP: srv.Client()}
	_, err := c.LookupWrappingToken("wraps.bad")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
