package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/vaultdiff/internal/vault"
)

func capabilitiesServer(t *testing.T, path string, caps []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/capabilities-self" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				path: caps,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestGetPolicyForPath_ReturnsCapabilities(t *testing.T) {
	path := "secret/data/myapp/config"
	caps := []string{"read", "list"}
	srv := capabilitiesServer(t, path, caps)
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{Address: srv.URL, Mount: "secret"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	accesses, err := client.GetPolicyForPath(context.Background(), path)
	if err != nil {
		t.Fatalf("GetPolicyForPath: %v", err)
	}
	if len(accesses) == 0 {
		t.Fatal("expected at least one access entry")
	}
	if accesses[0].Path != path {
		t.Errorf("expected path %q, got %q", path, accesses[0].Path)
	}
}

func TestCanRead_TrueForReadCap(t *testing.T) {
	if !vault.CanRead([]string{"read", "list"}) {
		t.Error("expected CanRead to return true")
	}
}

func TestCanRead_FalseForDeny(t *testing.T) {
	if vault.CanRead([]string{"deny"}) {
		t.Error("expected CanRead to return false for deny")
	}
}

func TestCanWrite_TrueForUpdate(t *testing.T) {
	if !vault.CanWrite([]string{"update"}) {
		t.Error("expected CanWrite to return true for update")
	}
}

func TestCanWrite_TrueForSudo(t *testing.T) {
	if !vault.CanWrite([]string{"sudo"}) {
		t.Error("expected CanWrite to return true for sudo")
	}
}

func TestCanWrite_FalseForReadOnly(t *testing.T) {
	if vault.CanWrite([]string{"read", "list"}) {
		t.Error("expected CanWrite to return false for read-only caps")
	}
}
