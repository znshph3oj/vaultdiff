package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func replicationServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetReplicationStatus_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"dr": map[string]any{
				"mode":              "primary",
				"primary":           true,
				"known_secondaries": []string{"sec1", "sec2"},
			},
			"performance": map[string]any{
				"mode": "secondary",
			},
		},
	}
	srv := replicationServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	rs, err := GetReplicationStatus(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rs.DRMode != "primary" {
		t.Errorf("expected dr_mode=primary, got %s", rs.DRMode)
	}
	if rs.PerformanceMode != "secondary" {
		t.Errorf("expected performance_mode=secondary, got %s", rs.PerformanceMode)
	}
	if !rs.Primary {
		t.Error("expected primary=true")
	}
	if len(rs.KnownSecondaries) != 2 {
		t.Errorf("expected 2 secondaries, got %d", len(rs.KnownSecondaries))
	}
}

func TestGetReplicationStatus_NotFound(t *testing.T) {
	srv := replicationServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetReplicationStatus(client)
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestGetReplicationStatus_UnexpectedStatus(t *testing.T) {
	srv := replicationServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test", HTTP: srv.Client()}
	_, err := GetReplicationStatus(client)
	if err == nil {
		t.Fatal("expected error for 500")
	}
}
