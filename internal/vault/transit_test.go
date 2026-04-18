package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func transitServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetTransitKeyInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type":                   "aes256-gcm96",
			"deletion_allowed":        false,
			"exportable":              true,
			"latest_version":          3,
			"min_decryption_version":  1,
		},
	}
	srv := transitServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	info, err := GetTransitKeyInfo(client, "mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Type != "aes256-gcm96" {
		t.Errorf("expected type aes256-gcm96, got %s", info.Type)
	}
	if info.LatestVersion != 3 {
		t.Errorf("expected latest_version 3, got %d", info.LatestVersion)
	}
	if !info.Exportable {
		t.Error("expected exportable true")
	}
}

func TestGetTransitKeyInfo_NotFound(t *testing.T) {
	srv := transitServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetTransitKeyInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetTransitKeyInfo_UnexpectedStatus(t *testing.T) {
	srv := transitServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetTransitKeyInfo(client, "mykey")
	if err == nil {
		t.Fatal("expected error for 500")
	}
}
