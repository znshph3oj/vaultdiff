package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func policyListServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/policies/acl":
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"keys": []string{"default", "admin"}},
			})
		case "/v1/sys/policies/acl/admin":
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"name": "admin", "rules": "path \"*\" { capabilities = [\"sudo\"] }"},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestListPolicies_Success(t *testing.T) {
	srv := policyListServer(t)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", "")
	policies, err := c.ListPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(policies))
	}
}

func TestGetPolicy_Success(t *testing.T) {
	srv := policyListServer(t)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", "")
	p, err := c.GetPolicy("admin")
	if err != nil {
		t.Fatal(err)
	}
	if p == nil || p.Name != "admin" {
		t.Fatalf("expected admin policy, got %v", p)
	}
}

func TestGetPolicy_NotFound(t *testing.T) {
	srv := policyListServer(t)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", "")
	p, err := c.GetPolicy("missing")
	if err != nil {
		t.Fatal(err)
	}
	if p != nil {
		t.Fatal("expected nil for missing policy")
	}
}
