package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func rabbitmqServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetRabbitMQRoleInfo_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"vhost": "/",
			"tags":  "administrator",
		},
	}
	srv := rabbitmqServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	info, err := GetRabbitMQRoleInfo(client, "rabbitmq", "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Vhost != "/" {
		t.Errorf("expected vhost '/', got %q", info.Vhost)
	}
	if info.Tags != "administrator" {
		t.Errorf("expected tags 'administrator', got %q", info.Tags)
	}
	if info.Name != "my-role" {
		t.Errorf("expected name 'my-role', got %q", info.Name)
	}
}

func TestGetRabbitMQRoleInfo_NotFound(t *testing.T) {
	srv := rabbitmqServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetRabbitMQRoleInfo(client, "rabbitmq", "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRabbitMQRoleInfo_UnexpectedStatus(t *testing.T) {
	srv := rabbitmqServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := GetRabbitMQRoleInfo(client, "rabbitmq", "my-role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
