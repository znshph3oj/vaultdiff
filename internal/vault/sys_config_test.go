package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sysConfigServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetSysConfig_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"default_lease_ttl": "768h",
			"max_lease_ttl":     "768h",
			"force_no_cache":    false,
		},
	}
	srv := sysConfigServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	cfg, err := c.GetSysConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultLeaseTTL != "768h" {
		t.Errorf("expected 768h, got %s", cfg.DefaultLeaseTTL)
	}
}

func TestGetSysConfig_NotFound(t *testing.T) {
	srv := sysConfigServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := c.GetSysConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSysConfig_UnexpectedStatus(t *testing.T) {
	srv := sysConfigServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := c.GetSysConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
