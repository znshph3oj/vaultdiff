package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func tagsServer(t *testing.T, customMeta map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			payload := map[string]interface{}{
				"data": map[string]interface{}{
					"custom_metadata": customMeta,
				},
			}
			_ = json.NewEncoder(w).Encode(payload)
		case r.Method == http.MethodPost || r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestGetSecretTags_ReturnsTags(t *testing.T) {
	srv := tagsServer(t, map[string]interface{}{"env": "prod", "team": "platform"})
	defer srv.Close()

	c, err := NewClient(srv.URL, "token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tags, err := c.GetSecretTags("myapp/config")
	if err != nil {
		t.Fatalf("GetSecretTags: %v", err)
	}
	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", tags["env"])
	}
	if tags["team"] != "platform" {
		t.Errorf("expected team=platform, got %q", tags["team"])
	}
}

func TestGetSecretTags_EmptyCustomMetadata(t *testing.T) {
	srv := tagsServer(t, nil)
	defer srv.Close()

	c, err := NewClient(srv.URL, "token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tags, err := c.GetSecretTags("myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestGetSecretTags_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c, err := NewClient(srv.URL, "token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.GetSecretTags("missing/secret")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
