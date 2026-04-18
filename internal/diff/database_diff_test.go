package diff

import (
	"bytes"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

var baseDatabaseRole = &vault.DatabaseRoleInfo{
	DBName:     "postgres",
	DefaultTTL: 3600,
	MaxTTL:     7200,
}

func TestCompareDatabaseRoles_NoChanges(t *testing.T) {
	changes := CompareDatabaseRoles(baseDatabaseRole, baseDatabaseRole)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareDatabaseRoles_DBNameChanged(t *testing.T) {
	b := *baseDatabaseRole
	b.DBName = "mysql"
	changes := CompareDatabaseRoles(baseDatabaseRole, &b)
	if len(changes) != 1 || changes[0].Field != "db_name" {
		t.Errorf("expected db_name change, got %+v", changes)
	}
}

func TestCompareDatabaseRoles_TTLChanged(t *testing.T) {
	b := *baseDatabaseRole
	b.DefaultTTL = 1800
	changes := CompareDatabaseRoles(baseDatabaseRole, &b)
	if len(changes) != 1 || changes[0].Field != "default_ttl" {
		t.Errorf("expected default_ttl change, got %+v", changes)
	}
}

func TestCompareDatabaseRoles_NilInputs(t *testing.T) {
	if CompareDatabaseRoles(nil, baseDatabaseRole) != nil {
		t.Error("expected nil for nil input a")
	}
	if CompareDatabaseRoles(baseDatabaseRole, nil) != nil {
		t.Error("expected nil for nil input b")
	}
}

func TestFprintDatabaseDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintDatabaseDiff(&buf, nil)
	if buf.String() == "" {
		t.Error("expected output for no changes")
	}
}

func TestFprintDatabaseDiff_WithChanges(t *testing.T) {
	changes := []DatabaseRoleChange{{Field: "db_name", OldValue: "postgres", NewValue: "mysql"}}
	var buf bytes.Buffer
	FprintDatabaseDiff(&buf, changes)
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
