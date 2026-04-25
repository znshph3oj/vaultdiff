package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func aliasServer(t *testing.T, aliasID string, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/identity/entity-alias/id/" + aliasID
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetAliasInfo_Success(t *testing.T) {
	aliasID := "abc-123"
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":             aliasID,
			"name":           "testuser",
			"mount_accessor": "auth_userpass_abc",
			"mount_type":     "userpass",
			"canonical_id":   "entity-456",
			"metadata":       map[string]string{"env": "prod"},
		},
	}
	svr := aliasServer(t, aliasID, http.StatusOK, payload)
	defer svr.Close()

	client, err := NewClient(svr.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := GetAliasInfo(client, aliasID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "testuser" {
		t.Errorf("expected name %q, got %q", "testuser", info.Name)
	}
	if info.MountType != "userpass" {
		t.Errorf("expected mount_type %q, got %q", "userpass", info.MountType)
	}
	if info.Metadata["env"] != "prod" {
		t.Errorf("expected metadata env=prod, got %q", info.Metadata["env"])
	}
}

func TestGetAliasInfo_NotFound(t *testing.T) {
	svr := aliasServer(t, "missing", http.StatusNotFound, nil)
	defer svr.Close()

	client, _ := NewClient(svr.URL, "test-token", "")
	_, err := GetAliasInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAliasInfo_UnexpectedStatus(t *testing.T) {
	svr := aliasServer(t, "bad", http.StatusInternalServerError, nil)
	defer svr.Close()

	client, _ := NewClient(svr.URL, "test-token", "")
	_, err := GetAliasInfo(client, "bad")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
