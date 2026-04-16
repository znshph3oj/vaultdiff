package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func pluginServer(t *testing.T, status int, plugins []PluginInfo) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status == http.StatusOK {
			body := map[string]interface{}{
				"data": map[string]interface{}{
					"detailed": plugins,
				},
			}
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListPlugins_Success(t *testing.T) {
	plugins := []PluginInfo{
		{Name: "aws", Type: "secret", Version: "v1.0.0", Builtin: true},
		{Name: "jwt", Type: "auth", Version: "v0.9.0", Builtin: false},
	}
	srv := pluginServer(t, http.StatusOK, plugins)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	result, err := c.ListPlugins("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(result))
	}
	if result[0].Name != "aws" {
		t.Errorf("expected aws, got %s", result[0].Name)
	}
}

func TestListPlugins_NotFound(t *testing.T) {
	srv := pluginServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	_, err := c.ListPlugins("secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListPlugins_UnexpectedStatus(t *testing.T) {
	srv := pluginServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	_, err := c.ListPlugins("auth")
	if err == nil {
		t.Fatal("expected error for forbidden status")
	}
}
