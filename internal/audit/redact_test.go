package audit

import (
	"testing"
	"time"
)

func baseEntry() Entry {
	return Entry{
		Timestamp: time.Now(),
		User:      "alice",
		Operation: "read",
		Path:      "secret/data/app",
		Data: map[string]string{
			"password": "s3cr3t",
			"username": "admin",
			"api_key":  "tok_abc123",
		},
	}
}

func TestRedact_ByKey(t *testing.T) {
	e := baseEntry()
	out := Redact(e, RedactOptions{Keys: []string{"password", "api_key"}})

	if out.Data["password"] != "[REDACTED]" {
		t.Errorf("expected password to be redacted, got %s", out.Data["password"])
	}
	if out.Data["api_key"] != "[REDACTED]" {
		t.Errorf("expected api_key to be redacted, got %s", out.Data["api_key"])
	}
	if out.Data["username"] == "[REDACTED]" {
		t.Error("username should not be redacted")
	}
}

func TestRedact_ByPattern(t *testing.T) {
	e := baseEntry()
	out := Redact(e, RedactOptions{Patterns: []string{`^tok_`}})

	if out.Data["api_key"] != "[REDACTED]" {
		t.Errorf("expected api_key to be redacted by pattern, got %s", out.Data["api_key"])
	}
	if out.Data["password"] == "[REDACTED]" {
		t.Error("password should not be redacted by pattern")
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	e := baseEntry()
	out := Redact(e, RedactOptions{Keys: []string{"PASSWORD"}})

	if out.Data["password"] != "[REDACTED]" {
		t.Errorf("expected case-insensitive match for PASSWORD, got %s", out.Data["password"])
	}
}

func TestRedact_NilData(t *testing.T) {
	e := Entry{User: "bob", Path: "secret/data/test"}
	out := Redact(e, RedactOptions{Keys: []string{"password"}})
	if out.Data != nil {
		t.Error("expected nil Data to remain nil")
	}
}

func TestRedactAll_AppliestoAll(t *testing.T) {
	entries := []Entry{baseEntry(), baseEntry()}
	out := RedactAll(entries, RedactOptions{Keys: []string{"password"}})

	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for i, e := range out {
		if e.Data["password"] != "[REDACTED]" {
			t.Errorf("entry %d: expected password redacted", i)
		}
	}
}

func TestRedact_OriginalUnmodified(t *testing.T) {
	e := baseEntry()
	_ = Redact(e, RedactOptions{Keys: []string{"password"}})

	if e.Data["password"] == "[REDACTED]" {
		t.Error("original entry should not be modified")
	}
}
