package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func databaseServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetDatabaseRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"db_name":             "mydb",
			"creation_statements": []string{"CREATE USER ..."},
			"default_ttl":         "1h",
			"max_ttl":             "24h",
		},
	}
	srv := databaseServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	info, err := GetDatabaseRoleInfo(client, "readonly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.DBName != "mydb" {
		t.Errorf("expected db_name=mydb, got %s", info.DBName)
	}
	if info.DefaultTTL != "1h" {
		t.Errorf("expected default_ttl=1h, got %s", info.DefaultTTL)
	}
	if info.MaxTTL != "24h" {
		t.Errorf("expected max_ttl=24h, got %s", info.MaxTTL)
	}
	if info.Name != "readonly" {
		t.Errorf("expected name=readonly, got %s", info.Name)
	}
}

func TestGetDatabaseRoleInfo_NotFound(t *testing.T) {
	srv := databaseServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetDatabaseRoleInfo(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetDatabaseRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := databaseServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetDatabaseRoleInfo(client, "role")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
