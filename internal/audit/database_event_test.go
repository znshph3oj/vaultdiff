package audit

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

func TestRecordDatabaseRoleChange_WritesEntry(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	changes := []diff.DatabaseRoleChange{
		{Field: "db_name", OldValue: "postgres", NewValue: "mysql"},
	}
	if err := RecordDatabaseRoleChange(logger, "database/roles/myrole", "admin", changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entry Entry
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("failed to decode entry: %v", err)
	}
	if entry.Operation != "database_role_diff" {
		t.Errorf("expected operation database_role_diff, got %s", entry.Operation)
	}
	if entry.Path != "database/roles/myrole" {
		t.Errorf("unexpected path: %s", entry.Path)
	}
	if entry.Data["db_name"] == "" {
		t.Error("expected db_name in data")
	}
}

func TestRecordDatabaseRoleChange_NoChanges_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	if err := RecordDatabaseRoleChange(logger, "database/roles/myrole", "admin", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected no output for empty changes")
	}
}
