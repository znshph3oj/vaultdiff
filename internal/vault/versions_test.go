package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListVersions_Success(t *testing.T) {
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"versions": map[string]interface{}{
				"1": map[string]interface{}{
					"created_time":  "2024-01-01T10:00:00Z",
					"deletion_time": "",
					"destroyed":     false,
				},
				"2": map[string]interface{}{
					"created_time":  "2024-01-02T10:00:00Z",
					"deletion_time": "",
					"destroyed":     false,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	versions, err := client.ListVersions("myapp/config")
	if err != nil {
		t.Fatalf("ListVersions error: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}

	// Versions should be sorted by created_time ascending.
	if versions[0].CreatedTime >= versions[1].CreatedTime {
		t.Errorf("expected versions sorted by created_time ascending")
	}

	if versions[0].Version != 1 || versions[1].Version != 2 {
		t.Errorf("expected sequential version numbers, got %d and %d",
			versions[0].Version, versions[1].Version)
	}
}

func TestListVersions_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.ListVersions("nonexistent/path")
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}
