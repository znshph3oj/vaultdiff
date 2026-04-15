package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/vaultdiff/internal/diff"
)

func alwaysRead(caps []string) bool  { return true }
func alwaysWrite(caps []string) bool { return false }

func TestNewAccessReport_SetsUser(t *testing.T) {
	r := diff.NewAccessReport("alice", map[string][]string{}, alwaysRead, alwaysWrite)
	if r.User != "alice" {
		t.Errorf("expected user alice, got %q", r.User)
	}
}

func TestNewAccessReport_PopulatesEntries(t *testing.T) {
	pathCaps := map[string][]string{
		"secret/data/app": {"read", "list"},
	}
	r := diff.NewAccessReport("bob", pathCaps, alwaysRead, alwaysWrite)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	entry := r.Entries[0]
	if entry.Path != "secret/data/app" {
		t.Errorf("unexpected path: %q", entry.Path)
	}
	if !entry.CanRead {
		t.Error("expected CanRead true")
	}
	if entry.CanWrite {
		t.Error("expected CanWrite false")
	}
}

func TestFprintAccessReport_ContainsHeader(t *testing.T) {
	pathCaps := map[string][]string{
		"secret/data/cfg": {"read"},
	}
	r := diff.NewAccessReport("carol", pathCaps, alwaysRead, alwaysWrite)
	var buf bytes.Buffer
	diff.FprintAccessReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "carol") {
		t.Error("output should contain user name")
	}
	if !strings.Contains(out, "PATH") {
		t.Error("output should contain PATH header")
	}
	if !strings.Contains(out, "secret/data/cfg") {
		t.Error("output should contain the path")
	}
}

func TestFprintAccessReport_MultipleEntries(t *testing.T) {
	pathCaps := map[string][]string{
		"secret/data/a": {"read"},
		"secret/data/b": {"create", "update"},
	}
	r := diff.NewAccessReport("dave", pathCaps, alwaysRead, alwaysWrite)
	var buf bytes.Buffer
	diff.FprintAccessReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "secret/data/a") {
		t.Error("expected secret/data/a in output")
	}
	if !strings.Contains(out, "secret/data/b") {
		t.Error("expected secret/data/b in output")
	}
}
