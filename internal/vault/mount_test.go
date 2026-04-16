package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mountServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/mounts" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListMounts_Success(t *testing.T) {
	payload := map[string]interface{}{
		"secret/": map[string]interface{}{
			"type":        "kv",
			"description": "key/value store",
			"options":     map[string]string{"version": "2"},
			"local":       false,
			"seal_wrap":   false,
		},
		"pki/": map[string]interface{}{
			"type":        "pki",
			"description": "PKI engine",
			"local":       true,
		},
	}
	srv := mountServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test-token", http: srv.Client()}
	mounts, err := c.ListMounts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mounts) != 2 {
		t.Fatalf("expected 2 mounts, got %d", len(mounts))
	}
	if mounts["secret/"].Type != "kv" {
		t.Errorf("expected type kv, got %s", mounts["secret/"].Type)
	}
	if mounts["pki/"].Local != true {
		t.Errorf("expected pki mount to be local")
	}
}

func TestListMounts_NotFound(t *testing.T) {
	srv := mountServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test-token", http: srv.Client()}
	_, err := c.ListMounts()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListMounts_UnexpectedStatus(t *testing.T) {
	srv := mountServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test-token", http: srv.Client()}
	_, err := c.ListMounts()
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
