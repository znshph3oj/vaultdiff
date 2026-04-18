package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func groupServer(t *testing.T, groupID string, status int, group *IdentityGroup) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/identity/group/id/" + groupID
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if group != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": group})
		}
	}))
}

func TestGetIdentityGroup_Success(t *testing.T) {
	expected := &IdentityGroup{
		ID:       "abc-123",
		Name:     "dev-team",
		Type:     "internal",
		Policies: []string{"default", "dev"},
		MemberEntityIDs: []string{"eid-1", "eid-2"},
	}
	srv := groupServer(t, "abc-123", http.StatusOK, expected)
	defer srv.Close()

	client := testClient(t, srv.URL)
	got, err := GetIdentityGroup(client, "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != expected.Name {
		t.Errorf("expected name %q, got %q", expected.Name, got.Name)
	}
	if len(got.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(got.Policies))
	}
}

func TestGetIdentityGroup_NotFound(t *testing.T) {
	srv := groupServer(t, "missing", http.StatusNotFound, nil)
	defer srv.Close()

	client := testClient(t, srv.URL)
	_, err := GetIdentityGroup(client, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetIdentityGroup_UnexpectedStatus(t *testing.T) {
	srv := groupServer(t, "abc-123", http.StatusInternalServerError, nil)
	defer srv.Close()

	client := testClient(t, srv.URL)
	_, err := GetIdentityGroup(client, "abc-123")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}
