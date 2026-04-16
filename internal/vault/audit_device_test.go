package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func auditDeviceServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListAuditDevices_Success(t *testing.T) {
	payload := map[string]interface{}{
		"file/": map[string]interface{}{
			"path":        "file/",
			"type":        "file",
			"description": "file audit log",
			"options":     map[string]string{"file_path": "/var/log/vault.log"},
			"local":       false,
		},
	}
	srv := auditDeviceServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	devices, err := c.ListAuditDevices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	dev, ok := devices["file/"]
	if !ok {
		t.Fatal("expected key 'file/'")
	}
	if dev.Type != "file" {
		t.Errorf("expected type 'file', got %q", dev.Type)
	}
}

func TestListAuditDevices_NotFound(t *testing.T) {
	srv := auditDeviceServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	_, err := c.ListAuditDevices()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListAuditDevices_UnexpectedStatus(t *testing.T) {
	srv := auditDeviceServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token")
	_, err := c.ListAuditDevices()
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
