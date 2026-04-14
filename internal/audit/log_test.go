package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{
		Timestamp:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Path:        "secret/data/myapp",
		FromVersion: 2,
		ToVersion:   3,
		Changes: []diff.Change{
			{Key: "DB_PASS", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		},
		User: "alice",
	}

	if err := l.Record(e); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	if !strings.HasSuffix(line, "}") {
		t.Fatalf("expected JSON object, got: %s", line)
	}

	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Path != e.Path {
		t.Errorf("Path: got %q, want %q", got.Path, e.Path)
	}
	if got.User != e.User {
		t.Errorf("User: got %q, want %q", got.User, e.User)
	}
	if len(got.Changes) != 1 {
		t.Fatalf("Changes: got %d, want 1", len(got.Changes))
	}
	if got.Changes[0].Key != "DB_PASS" {
		t.Errorf("Change key: got %q, want DB_PASS", got.Changes[0].Key)
	}
}

func TestRecord_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	before := time.Now().UTC()
	if err := l.Record(audit.Entry{Path: "secret/data/test"}); err != nil {
		t.Fatalf("Record() error: %v", err)
	}
	after := time.Now().UTC()

	var got audit.Entry
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("Timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestNewLogger_DefaultsToStdout(t *testing.T) {
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil Logger")
	}
}
