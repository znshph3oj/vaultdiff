package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewReport_SummarizesChanges(t *testing.T) {
	results := []DiffResult{
		{Key: "FOO", Status: StatusAdded, NewValue: "bar"},
		{Key: "BAZ", Status: StatusRemoved, OldValue: "old"},
		{Key: "QUX", Status: StatusModified, OldValue: "a", NewValue: "b"},
		{Key: "SAME", Status: StatusUnchanged, OldValue: "x", NewValue: "x"},
	}

	report := NewReport("secret/myapp", 1, 2, results)

	if report.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", report.Summary.Added)
	}
	if report.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", report.Summary.Removed)
	}
	if report.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", report.Summary.Modified)
	}
	if report.Summary.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", report.Summary.Unchanged)
	}
}

func TestNewReport_MetadataIsSet(t *testing.T) {
	report := NewReport("secret/myapp", 3, 5, nil)

	if report.Path != "secret/myapp" {
		t.Errorf("unexpected path: %s", report.Path)
	}
	if report.VersionA != 3 || report.VersionB != 5 {
		t.Errorf("unexpected versions: %d, %d", report.VersionA, report.VersionB)
	}
	if report.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestReport_WriteJSON(t *testing.T) {
	report := NewReport("secret/myapp", 1, 2, []DiffResult{
		{Key: "TOKEN", Status: StatusModified, OldValue: "old", NewValue: "new"},
	})

	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}
	if decoded.Path != "secret/myapp" {
		t.Errorf("unexpected decoded path: %s", decoded.Path)
	}
	if len(decoded.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(decoded.Changes))
	}
}

func TestReport_PrintSummary(t *testing.T) {
	report := NewReport("secret/myapp", 1, 2, []DiffResult{
		{Key: "A", Status: StatusAdded},
		{Key: "B", Status: StatusRemoved},
	})

	var buf bytes.Buffer
	report.PrintSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "secret/myapp") {
		t.Errorf("summary missing path, got: %s", out)
	}
	if !strings.Contains(out, "+1 added") {
		t.Errorf("summary missing added count, got: %s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("summary missing removed count, got: %s", out)
	}
}
