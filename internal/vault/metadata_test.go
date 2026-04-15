package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func metadataServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/"+path {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"data": payload})
			return
		}
		http.NotFound(w, r)
	}))
}

func TestGetSecretMetadata_Success(t *testing.T) {
	payload := map[string]interface{}{
		"current_version": float64(3),
		"oldest_version":  float64(1),
		"created_time":    "2024-01-01T00:00:00Z",
		"updated_time":    "2024-06-01T00:00:00Z",
		"versions": map[string]interface{}{
			"1": map[string]interface{}{
				"created_time": "2024-01-01T00:00:00Z",
				"deletion_time": "",
				"destroyed":    false,
			},
			"3": map[string]interface{}{
				"created_time": "2024-06-01T00:00:00Z",
				"deletion_time": "",
				"destroyed":    false,
			},
		},
	}

	srv := metadataServer(t, "secret/metadata/myapp/config", payload)
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	vc, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: vc.Logical(), mount: "secret"}

	meta, err := c.GetSecretMetadata(context.Background(), "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.CurrentVersion != 3 {
		t.Errorf("expected current_version 3, got %d", meta.CurrentVersion)
	}
	if meta.OldestVersion != 1 {
		t.Errorf("expected oldest_version 1, got %d", meta.OldestVersion)
	}
	if len(meta.Versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(meta.Versions))
	}
}

func TestGetSecretMetadata_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	vc, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: vc.Logical(), mount: "secret"}

	_, err := c.GetSecretMetadata(context.Background(), "missing/path")
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestGetSecretMetadata_VersionDestroyed(t *testing.T) {
	payload := map[string]interface{}{
		"current_version": float64(2),
		"oldest_version":  float64(1),
		"created_time":    "2024-01-01T00:00:00Z",
		"updated_time":    "2024-05-01T00:00:00Z",
		"versions": map[string]interface{}{
			"1": map[string]interface{}{
				"created_time":  "2024-01-01T00:00:00Z",
				"deletion_time": "2024-03-01T00:00:00Z",
				"destroyed":     true,
			},
		},
	}

	srv := metadataServer(t, "secret/metadata/myapp/creds", payload)
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	vc, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: vc.Logical(), mount: "secret"}

	meta, err := c.GetSecretMetadata(context.Background(), "myapp/creds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v1, ok := meta.Versions[1]
	if !ok {
		t.Fatal("expected version 1 to exist")
	}
	if !v1.Destroyed {
		t.Error("expected version 1 to be marked destroyed")
	}
}
