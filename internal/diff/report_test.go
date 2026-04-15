package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewReport_SummarizesChanges(t *testing.T) {
	changes := []DiffResult{
		{Key: "a", Status: StatusAdded},
		{Key: "b", Status: StatusRemoved},
		{Key: "c", Status: StatusModified},
		{Key: "d", Status: StatusUnchanged},
		{Key: "e", Status: StatusAdded},
	}
	r := NewReport("secret/data/app", 1, 2, changes)
	if r.Added != 2 {
		t.Errorf("expected Added=2, got %d", r.Added)
	}
	if r.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", r.Removed)
	}
	if r.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", r.Modified)
	}
	if r.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", r.Unchanged)
	}
}

func TestNewReport_MetadataIsSet(t *testing.T) {
	r := NewReport("secret/data/app", 3, 4, nil)
	if r.Path != "secret/data/app" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.FromVersion != 3 || r.ToVersion != 4 {
		t.Errorf("unexpected versions: %d %d", r.FromVersion, r.ToVersion)
	}
	if r.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should not be zero")
	}
}

func TestReport_WriteJSON(t *testing.T) {
	changes := []DiffResult{{Key: "token", Status: StatusModified, OldValue: "x", NewValue: "y"}}
	r := NewReport("secret/data/svc", 1, 2, changes)
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["path"] != "secret/data/svc" {
		t.Errorf("unexpected path in JSON: %v", out["path"])
	}
}

func TestReport_PrintSummary(t *testing.T) {
	r := NewReport("secret/data/db", 2, 3, []DiffResult{
		{Key: "pass", Status: StatusModified},
	})
	var buf bytes.Buffer
	r.PrintSummary(&buf)
	out := buf.String()
	if !strings.Contains(out, "secret/data/db") {
		t.Error("summary should contain path")
	}
	if !strings.Contains(out, "Modified: 1") {
		t.Error("summary should contain modified count")
	}
}
