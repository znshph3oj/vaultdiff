package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func awsServer(t *testing.T, role string, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/aws/roles/" + role
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetAWSRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"credential_type": "assumed_role",
			"role_arns":       []string{"arn:aws:iam::123456789012:role/MyRole"},
			"policy_arns":     []string{},
			"default_ttl":     3600,
			"max_ttl":         7200,
		},
	}
	srv := awsServer(t, "my-role", http.StatusOK, payload)
	defer srv.Close()

	client := testClient(t, srv.URL)
	info, err := GetAWSRoleInfo(client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.CredentialType != "assumed_role" {
		t.Errorf("expected assumed_role, got %s", info.CredentialType)
	}
	if info.Name != "my-role" {
		t.Errorf("expected name my-role, got %s", info.Name)
	}
	if info.DefaultTTL != 3600 {
		t.Errorf("expected default_ttl 3600, got %d", info.DefaultTTL)
	}
}

func TestGetAWSRoleInfo_NotFound(t *testing.T) {
	srv := awsServer(t, "missing", http.StatusNotFound, nil)
	defer srv.Close()

	client := testClient(t, srv.URL)
	_, err := GetAWSRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAWSRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := awsServer(t, "bad", http.StatusInternalServerError, nil)
	defer srv.Close()

	client := testClient(t, srv.URL)
	_, err := GetAWSRoleInfo(client, "bad")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
