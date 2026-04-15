package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSnapshot_SetsFields(t *testing.T) {
	data := map[string]string{"key": "value"}
	s := NewSnapshot("secret/myapp", 3, data)

	if s.Path != "secret/myapp" {
		t.Errorf("expected path 'secret/myapp', got %q", s.Path)
	}
	if s.Version != 3 {
		t.Errorf("expected version 3, got %d", s.Version)
	}
	if s.Data["key"] != "value" {
		t.Errorf("expected data key 'value', got %q", s.Data["key"])
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := &Snapshot{
		Path:       "secret/test",
		Version:    2,
		CapturedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Data:       map[string]string{"foo": "bar", "baz": "qux"},
	}

	if err := SaveSnapshot(orig, path); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if loaded.Path != orig.Path {
		t.Errorf("path mismatch: got %q, want %q", loaded.Path, orig.Path)
	}
	if loaded.Version != orig.Version {
		t.Errorf("version mismatch: got %d, want %d", loaded.Version, orig.Version)
	}
	if loaded.Data["foo"] != "bar" {
		t.Errorf("data mismatch: got %q, want 'bar'", loaded.Data["foo"])
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	s := NewSnapshot("secret/x", 1, map[string]string{})
	err := SaveSnapshot(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestDiffSnapshots_DetectsChanges(t *testing.T) {
	old := NewSnapshot("secret/app", 1, map[string]string{"a": "1", "b": "2"})
	new := NewSnapshot("secret/app", 2, map[string]string{"a": "1", "b": "changed", "c": "3"})

	changes := DiffSnapshots(old, new)

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}

	changeMap := map[string]ChangeType{}
	for _, c := range changes {
		changeMap[c.Key] = c.Type
	}

	if changeMap["b"] != Modified {
		t.Errorf("expected 'b' to be Modified")
	}
	if changeMap["c"] != Added {
		t.Errorf("expected 'c' to be Added")
	}

	_ = os.Remove("snap.json")
}
