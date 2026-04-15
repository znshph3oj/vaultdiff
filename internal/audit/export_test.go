package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var sampleExportEntries = []Entry{
	{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		User:      "alice",
		Path:      "secret/data/app",
		Operation: "diff",
		VersionA:  1,
		VersionB:  2,
	},
	{
		Timestamp: time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC),
		User:      "bob",
		Path:      "secret/data/db",
		Operation: "read",
		VersionA:  3,
		VersionB:  0,
	},
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, sampleExportEntries, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded []Entry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}

	if len(decoded) != 2 {
		t.Errorf("expected 2 entries, got %d", len(decoded))
	}
	if decoded[0].User != "alice" {
		t.Errorf("expected user alice, got %s", decoded[0].User)
	}
}

func TestExport_CSV(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, sampleExportEntries, FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 { // header + 2 rows
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "alice") {
		t.Errorf("expected alice in first data row, got: %s", lines[1])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, sampleExportEntries, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported export format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, []Entry{}, FormatCSV); err != nil {
		t.Fatalf("unexpected error for empty entries: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
