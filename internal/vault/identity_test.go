package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func identityServer(t *testing.T, entityID string, entity *IdentityEntity, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/identity/entity/id/" + entityID
		if r.URL.Path != expected {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
		if entity != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"data": entity})
		}
	}))
}

func TestGetIdentityEntity_Success(t *testing.T) {
	entity := &IdentityEntity{
		ID:       "abc-123",
		Name:     "alice",
		Policies: []string{"default", "admin"},
		Metadata: map[string]string{"team": "ops"},
		Disabled: false,
	}
	srv := identityServer(t, "abc-123", entity, http.StatusOK)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	got, err := c.GetIdentityEntity("abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "alice" {
		t.Errorf("expected name alice, got %s", got.Name)
	}
	if len(got.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(got.Policies))
	}
}

func TestGetIdentityEntity_NotFound(t *testing.T) {
	srv := identityServer(t, "abc-123", nil, http.StatusNotFound)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := c.GetIdentityEntity("abc-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetIdentityEntity_UnexpectedStatus(t *testing.T) {
	srv := identityServer(t, "abc-123", nil, http.StatusInternalServerError)
	defer srv.Close()

	c := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := c.GetIdentityEntity("abc-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
