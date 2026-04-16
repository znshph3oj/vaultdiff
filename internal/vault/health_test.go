package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func healthServer(status HealthStatus, code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(status)
	}))
}

func TestGetHealth_Success(t *testing.T) {
	srv := healthServer(HealthStatus{
		Initialized: true,
		Sealed:      false,
		Version:     "1.15.0",
		ClusterName: "vault-cluster",
	}, http.StatusOK)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test"}
	h, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !h.Initialized {
		t.Error("expected initialized to be true")
	}
	if h.Version != "1.15.0" {
		t.Errorf("expected version 1.15.0, got %s", h.Version)
	}
}

func TestGetHealth_Standby(t *testing.T) {
	srv := healthServer(HealthStatus{Initialized: true, Standby: true, Version: "1.15.0"}, http.StatusTooManyRequests)
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test"}
	h, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !h.Standby {
		t.Error("expected standby to be true")
	}
}

func TestGetHealth_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := &Client{address: srv.URL, token: "test"}
	_, err := c.GetHealth(context.Background())
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
