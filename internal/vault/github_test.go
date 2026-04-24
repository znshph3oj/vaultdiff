package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func githubServer(t *testing.T, team string, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/auth/github/map/teams/" + team
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetGitHubRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"team_name": "devops",
			"policies":  []string{"read", "write"},
			"ttl":       "1h",
			"max_ttl":   "24h",
		},
	}
	srv := githubServer(t, "devops", http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	info, err := GetGitHubRoleInfo(client, "github", "devops")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.TeamName != "devops" {
		t.Errorf("expected team_name=devops, got %q", info.TeamName)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.TTL != "1h" {
		t.Errorf("expected ttl=1h, got %q", info.TTL)
	}
}

func TestGetGitHubRoleInfo_NotFound(t *testing.T) {
	srv := githubServer(t, "devops", http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	_, err := GetGitHubRoleInfo(client, "github", "devops")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetGitHubRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := githubServer(t, "devops", http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token"}
	_, err := GetGitHubRoleInfo(client, "github", "devops")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
