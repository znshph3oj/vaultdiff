package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sshServer(t *testing.T, status int, payload interface{}) (*httptest.Server, *Client) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
	t.Cleanup(ts.Close)
	client := &Client{Address: ts.URL, Token: "test-token"}
	return ts, client
}

func TestGetSSHRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"key_type":     "ca",
			"default_user": "ubuntu",
			"ttl":          "1h",
			"max_ttl":      "24h",
		},
	}
	_, client := sshServer(t, http.StatusOK, payload)
	info, err := GetSSHRoleInfo(client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.KeyType != "ca" {
		t.Errorf("expected key_type ca, got %s", info.KeyType)
	}
	if info.DefaultUser != "ubuntu" {
		t.Errorf("expected default_user ubuntu, got %s", info.DefaultUser)
	}
}

func TestGetSSHRoleInfo_NotFound(t *testing.T) {
	_, client := sshServer(t, http.StatusNotFound, nil)
	_, err := GetSSHRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSSHRoleInfo_UnexpectedStatus(t *testing.T) {
	_, client := sshServer(t, http.StatusInternalServerError, nil)
	_, err := GetSSHRoleInfo(client, "broken")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
