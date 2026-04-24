package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func oktaServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetOktaRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"policies":     []string{"read", "write"},
			"ttl":          "1h",
			"max_ttl":      "24h",
			"bound_groups": []string{"eng"},
			"bound_users":  []string{"alice"},
		},
	}
	svr := oktaServer(t, http.StatusOK, payload)
	defer svr.Close()

	client := &Client{Address: svr.URL, Token: "test-token"}
	info, err := GetOktaRoleInfo(client, "okta", "eng-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.TTL != "1h" {
		t.Errorf("expected TTL 1h, got %s", info.TTL)
	}
	if len(info.BoundGroups) != 1 || info.BoundGroups[0] != "eng" {
		t.Errorf("unexpected bound_groups: %v", info.BoundGroups)
	}
}

func TestGetOktaRoleInfo_NotFound(t *testing.T) {
	svr := oktaServer(t, http.StatusNotFound, nil)
	defer svr.Close()

	client := &Client{Address: svr.URL, Token: "test-token"}
	_, err := GetOktaRoleInfo(client, "okta", "missing-role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetOktaRoleInfo_UnexpectedStatus(t *testing.T) {
	svr := oktaServer(t, http.StatusInternalServerError, nil)
	defer svr.Close()

	client := &Client{Address: svr.URL, Token: "test-token"}
	_, err := GetOktaRoleInfo(client, "okta", "some-role")
	if err == nil {
		t.Fatal("expected error for unexpected status")
	}
}
