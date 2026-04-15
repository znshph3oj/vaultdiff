package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func secretPayload(version int, data map[string]interface{}) []byte {
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"data": data,
			"metadata": map[string]interface{}{
				"version": float64(version),
			},
		},
	})
	return body
}

func TestGetSecretVersion_ReturnsData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/v1/secret/data/myapp/config")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(secretPayload(3, map[string]interface{}{"API_KEY": "abc123", "DB_PASS": "secret"}))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "fake-token", "")
	require.NoError(t, err)

	sv, err := client.GetSecretVersion(context.Background(), "myapp/config", 3)
	require.NoError(t, err)
	require.NotNil(t, sv)

	assert.Equal(t, "abc123", sv.Data["API_KEY"])
	assert.Equal(t, "secret", sv.Data["DB_PASS"])
	assert.Equal(t, 3, sv.Version)
}

func TestGetSecretVersion_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "fake-token", "")
	require.NoError(t, err)

	_, err = client.GetSecretVersion(context.Background(), "missing/path", 1)
	assert.Error(t, err)
}

func TestGetSecretVersion_LatestVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.Query().Get("version"), "expected no version param for latest")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(secretPayload(5, map[string]interface{}{"FOO": "bar"}))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "fake-token", "")
	require.NoError(t, err)

	sv, err := client.GetSecretVersion(context.Background(), "myapp/config", 0)
	require.NoError(t, err)
	assert.Equal(t, "bar", sv.Data["FOO"])
}
