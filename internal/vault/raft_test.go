package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func raftServer(status int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetRaftStatus_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"leader_id":     "node1",
			"applied_index": 42,
			"commit_index":  42,
			"servers": []map[string]interface{}{
				{"node_id": "node1", "address": "127.0.0.1:8201", "leader": true, "voter": true, "protocol_version": "3"},
			},
		},
	}
	srv := raftServer(http.StatusOK, payload)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test", http: srv.Client()}
	status, err := c.GetRaftStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.LeaderID != "node1" {
		t.Errorf("expected leader node1, got %s", status.LeaderID)
	}
	if len(status.Servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(status.Servers))
	}
}

func TestGetRaftStatus_NotFound(t *testing.T) {
	srv := raftServer(http.StatusNotFound, nil)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test", http: srv.Client()}
	_, err := c.GetRaftStatus()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRaftStatus_UnexpectedStatus(t *testing.T) {
	srv := raftServer(http.StatusInternalServerError, nil)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test", http: srv.Client()}
	_, err := c.GetRaftStatus()
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
