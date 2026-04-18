package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func kvEngineServer(t *testing.T, status int, options map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if status == http.StatusOK {
			json.NewEncoder(w).Encode(map[string]interface{}{"options": options})
		}
	}))
}

func TestGetKVEngineInfo_V2(t *testing.T) {
	srv := kvEngineServer(t, http.StatusOK, map[string]string{"version": "2"})
	defer srv.Close()
	c := &Client{Address: srv.URL, Token: "tok", HTTP: srv.Client()}
	info, err := GetKVEngineInfo(c, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if info.Version != 2 {
		t.Errorf("expected version 2, got %d", info.Version)
	}
}

func TestGetKVEngineInfo_V1Default(t *testing.T) {
	srv := kvEngineServer(t, http.StatusOK, map[string]string{})
	defer srv.Close()
	c := &Client{Address: srv.URL, Token: "tok", HTTP: srv.Client()}
	info, err := GetKVEngineInfo(c, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if info.Version != 1 {
		t.Errorf("expected version 1, got %d", info.Version)
	}
}

func TestGetKVEngineInfo_NotFound(t *testing.T) {
	srv := kvEngineServer(t, http.StatusNotFound, nil)
	defer srv.Close()
	c := &Client{Address: srv.URL, Token: "tok", HTTP: srv.Client()}
	_, err := GetKVEngineInfo(c, "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetKVEngineInfo_UnexpectedStatus(t *testing.T) {
	srv := kvEngineServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()
	c := &Client{Address: srv.URL, Token: "tok", HTTP: srv.Client()}
	_, err := GetKVEngineInfo(c, "secret")
	if err == nil {
		t.Fatal("expected error")
	}
}
