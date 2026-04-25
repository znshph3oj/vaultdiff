package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func terraformServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetTerraformRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"organization": "my-org",
			"team_id":      "team-abc",
			"ttl":          3600,
			"max_ttl":      7200,
		},
	}
	ts := terraformServer(t, http.StatusOK, payload)
	defer ts.Close()

	client, err := vault.NewClient(ts.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := vault.GetTerraformRoleInfo(client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Organization != "my-org" {
		t.Errorf("expected organization my-org, got %s", info.Organization)
	}
	if info.TeamID != "team-abc" {
		t.Errorf("expected team_id team-abc, got %s", info.TeamID)
	}
	if info.TTL != 3600 {
		t.Errorf("expected TTL 3600, got %d", info.TTL)
	}
	if info.MaxTTL != 7200 {
		t.Errorf("expected MaxTTL 7200, got %d", info.MaxTTL)
	}
}

func TestGetTerraformRoleInfo_NotFound(t *testing.T) {
	ts := terraformServer(t, http.StatusNotFound, nil)
	defer ts.Close()

	client, err := vault.NewClient(ts.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = vault.GetTerraformRoleInfo(client, "missing-role")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetTerraformRoleInfo_UnexpectedStatus(t *testing.T) {
	ts := terraformServer(t, http.StatusInternalServerError, nil)
	defer ts.Close()

	client, err := vault.NewClient(ts.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = vault.GetTerraformRoleInfo(client, "some-role")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
