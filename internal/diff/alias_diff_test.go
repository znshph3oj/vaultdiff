package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseAlias = &vault.AliasInfo{
	ID:            "alias-1",
	Name:          "alice",
	MountAccessor: "auth_userpass_xyz",
	MountType:     "userpass",
	CanonicalID:   "entity-99",
	Metadata:      map[string]string{"team": "platform"},
}

func TestCompareAliases_NoChanges(t *testing.T) {
	copy := *baseAlias
	copy.Metadata = map[string]string{"team": "platform"}
	changes := CompareAliases(baseAlias, &copy)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareAliases_NameChanged(t *testing.T) {
	updated := *baseAlias
	updated.Name = "bob"
	changes := CompareAliases(baseAlias, &updated)
	if len(changes) != 1 || changes[0].Field != "name" {
		t.Errorf("expected name change, got %+v", changes)
	}
}

func TestCompareAliases_MountTypeChanged(t *testing.T) {
	updated := *baseAlias
	updated.MountType = "github"
	changes := CompareAliases(baseAlias, &updated)
	if len(changes) != 1 || changes[0].Field != "mount_type" {
		t.Errorf("expected mount_type change, got %+v", changes)
	}
}

func TestCompareAliases_MetadataChanged(t *testing.T) {
	updated := *baseAlias
	updated.Metadata = map[string]string{"team": "security"}
	changes := CompareAliases(baseAlias, &updated)
	if len(changes) != 1 || changes[0].Field != "metadata.team" {
		t.Errorf("expected metadata.team change, got %+v", changes)
	}
}

func TestCompareAliases_NilInputs(t *testing.T) {
	if changes := CompareAliases(nil, nil); len(changes) != 0 {
		t.Errorf("expected no changes for nil inputs, got %d", len(changes))
	}
	changes := CompareAliases(nil, baseAlias)
	if len(changes) == 0 {
		t.Error("expected changes when comparing nil to non-nil")
	}
}

func TestFprintAliasDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintAliasDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No alias changes") {
		t.Errorf("expected no-changes message, got: %q", buf.String())
	}
}

func TestFprintAliasDiff_WithChanges(t *testing.T) {
	changes := []AliasDiffEntry{
		{Field: "name", Old: "alice", New: "bob"},
	}
	var buf bytes.Buffer
	FprintAliasDiff(&buf, changes)
	if !strings.Contains(buf.String(), "Alias Changes") {
		t.Errorf("expected header, got: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "alice") || !strings.Contains(buf.String(), "bob") {
		t.Errorf("expected old/new values in output, got: %q", buf.String())
	}
}
