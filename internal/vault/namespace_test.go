package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func namespaceServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if status == http.StatusOK {
			body := map[string]interface{}{
				"data": map[string]interface{}{"keys": keys},
			}
			json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListNamespaces_Success(t *testing.T) {
	srv := namespaceServer(t, http.StatusOK, []string{"team-a/", "team-b/"})
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	ns, err := c.ListNamespaces("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ns) != 2 {
		t.Fatalf("expected 2 namespaces, got %d", len(ns))
	}
	if ns[0].Path != "team-a/" {
		t.Errorf("expected team-a/, got %s", ns[0].Path)
	}
}

func TestListNamespaces_NotFound(t *testing.T) {
	srv := namespaceServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	_, err := c.ListNamespaces("missing/")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListNamespaces_UnexpectedStatus(t *testing.T) {
	srv := namespaceServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", "")
	_, err := c.ListNamespaces("")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
