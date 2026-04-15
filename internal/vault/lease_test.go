package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func leaseServer(t *testing.T, leaseID string, ttl int, expireTime string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"id":          leaseID,
				"renewable":   true,
				"ttl":         ttl,
				"expire_time": expireTime,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestGetLeaseInfo_Success(t *testing.T) {
	expire := time.Now().Add(5 * time.Minute).UTC().Format(time.RFC3339)
	srv := leaseServer(t, "database/creds/my-role/abc123", 300, expire, http.StatusOK)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := c.GetLeaseInfo("database/creds/my-role/abc123")
	if err != nil {
		t.Fatalf("GetLeaseInfo: %v", err)
	}

	if info.LeaseID != "database/creds/my-role/abc123" {
		t.Errorf("expected lease ID, got %q", info.LeaseID)
	}
	if !info.Renewable {
		t.Error("expected renewable to be true")
	}
	if info.LeaseDuration != 300*time.Second {
		t.Errorf("expected 300s duration, got %v", info.LeaseDuration)
	}
	if info.ExpireTime.IsZero() {
		t.Error("expected non-zero expire time")
	}
}

func TestGetLeaseInfo_NotFound(t *testing.T) {
	srv := leaseServer(t, "", 0, "", http.StatusNotFound)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.GetLeaseInfo("nonexistent/lease")
	if err == nil {
		t.Fatal("expected error for not found lease")
	}
}

func TestGetLeaseInfo_BadExpireTime(t *testing.T) {
	srv := leaseServer(t, "some/lease/id", 60, "not-a-date", http.StatusOK)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := c.GetLeaseInfo("some/lease/id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.ExpireTime.IsZero() {
		t.Error("expected zero expire time for invalid date string")
	}
}
