package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func quotaServer(t *testing.T, name string, status int, info *QuotaInfo) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/sys/quotas/rate-limit/" + name
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if info != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"data": info})
		}
	}))
}

func TestGetQuota_Success(t *testing.T) {
	info := &QuotaInfo{Name: "global", Path: "secret/", Type: "rate-limit", Rate: 100, Interval: 1}
	srv := quotaServer(t, "global", http.StatusOK, info)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	got, err := c.GetQuota("global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "global" {
		t.Errorf("expected name global, got %s", got.Name)
	}
	if got.Rate != 100 {
		t.Errorf("expected rate 100, got %f", got.Rate)
	}
}

func TestGetQuota_NotFound(t *testing.T) {
	srv := quotaServer(t, "missing", http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	_, err := c.GetQuota("missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetQuota_UnexpectedStatus(t *testing.T) {
	srv := quotaServer(t, "bad", http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	_, err := c.GetQuota("bad")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
