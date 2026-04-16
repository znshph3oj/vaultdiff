package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sealServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/seal-status" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetSealStatus_Unsealed(t *testing.T) {
	payload := SealStatus{
		Sealed:      false,
		Initialized: true,
		T:           1,
		N:           1,
		Version:     "1.13.0",
		ClusterName: "vault-cluster",
	}
	srv := sealServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	got, err := c.GetSealStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Sealed {
		t.Error("expected unsealed")
	}
	if got.Version != "1.13.0" {
		t.Errorf("expected version 1.13.0, got %s", got.Version)
	}
}

func TestGetSealStatus_Sealed(t *testing.T) {
	payload := SealStatus{Sealed: true, Progress: 1, T: 3, N: 5}
	srv := sealServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	got, err := c.GetSealStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Sealed {
		t.Error("expected sealed")
	}
}

func TestGetSealStatus_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	_, err := c.GetSealStatus()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSealStatus_UnexpectedStatus(t *testing.T) {
	srv := sealServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	_, err := c.GetSealStatus()
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
