package audit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/audit"
)

const sampleLog = `{"timestamp":"2024-01-10T08:00:00Z","path":"secret/data/app","from_version":1,"to_version":2,"changes":null,"user":"alice"}
{"timestamp":"2024-01-11T09:00:00Z","path":"secret/data/db","from_version":3,"to_version":4,"changes":null,"user":"bob"}
{"timestamp":"2024-01-12T10:00:00Z","path":"secret/data/app","from_version":2,"to_version":3,"changes":null,"user":"alice"}
`

func TestFilter_ByPath(t *testing.T) {
	entries, err := audit.Filter(strings.NewReader(sampleLog), audit.FilterOptions{Path: "secret/data/app"})
	if err != nil {
		t.Fatalf("Filter() error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("got %d entries, want 2", len(entries))
	}
}

func TestFilter_ByUser(t *testing.T) {
	entries, err := audit.Filter(strings.NewReader(sampleLog), audit.FilterOptions{User: "bob"})
	if err != nil {
		t.Fatalf("Filter() error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("got %d entries, want 1", len(entries))
	}
	if entries[0].User != "bob" {
		t.Errorf("expected user bob, got %q", entries[0].User)
	}
}

func TestFilter_BySince(t *testing.T) {
	since := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	entries, err := audit.Filter(strings.NewReader(sampleLog), audit.FilterOptions{Since: since})
	if err != nil {
		t.Fatalf("Filter() error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("got %d entries, want 2", len(entries))
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	entries, err := audit.Filter(strings.NewReader(sampleLog), audit.FilterOptions{})
	if err != nil {
		t.Fatalf("Filter() error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("got %d entries, want 3", len(entries))
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	entries, err := audit.Filter(strings.NewReader(""), audit.FilterOptions{})
	if err != nil {
		t.Fatalf("Filter() error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("got %d entries, want 0", len(entries))
	}
}
