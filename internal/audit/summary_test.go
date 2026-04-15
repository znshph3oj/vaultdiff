package audit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleEntries() []Entry {
	return []Entry{
		{Timestamp: time.Now(), User: "alice", Path: "secret/db", Operation: "diff"},
		{Timestamp: time.Now(), User: "bob", Path: "secret/app", Operation: "diff"},
		{Timestamp: time.Now(), User: "alice", Path: "secret/db", Operation: "watch"},
		{Timestamp: time.Now(), User: "alice", Path: "secret/app", Operation: "diff"},
	}
}

func TestSummarize_TotalCount(t *testing.T) {
	entries := sampleEntries()
	s := Summarize(entries)
	if s.TotalEvents != 4 {
		t.Errorf("expected 4 total events, got %d", s.TotalEvents)
	}
}

func TestSummarize_ByUser(t *testing.T) {
	s := Summarize(sampleEntries())
	if s.ByUser["alice"] != 3 {
		t.Errorf("expected alice=3, got %d", s.ByUser["alice"])
	}
	if s.ByUser["bob"] != 1 {
		t.Errorf("expected bob=1, got %d", s.ByUser["bob"])
	}
}

func TestSummarize_ByPath(t *testing.T) {
	s := Summarize(sampleEntries())
	if s.ByPath["secret/db"] != 2 {
		t.Errorf("expected secret/db=2, got %d", s.ByPath["secret/db"])
	}
	if s.ByPath["secret/app"] != 2 {
		t.Errorf("expected secret/app=2, got %d", s.ByPath["secret/app"])
	}
}

func TestSummarize_ByOperation(t *testing.T) {
	s := Summarize(sampleEntries())
	if s.ByOperation["diff"] != 3 {
		t.Errorf("expected diff=3, got %d", s.ByOperation["diff"])
	}
	if s.ByOperation["watch"] != 1 {
		t.Errorf("expected watch=1, got %d", s.ByOperation["watch"])
	}
}

func TestSummarize_EmptyInput(t *testing.T) {
	s := Summarize([]Entry{})
	if s.TotalEvents != 0 {
		t.Errorf("expected 0 events, got %d", s.TotalEvents)
	}
}

func TestPrintSummary_ContainsExpectedText(t *testing.T) {
	s := Summarize(sampleEntries())
	var buf bytes.Buffer
	PrintSummary(s, &buf)
	out := buf.String()
	for _, want := range []string{"Total events", "alice", "secret/db", "diff", "watch"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrintSummary_NilWriterUsesStdout(t *testing.T) {
	// Just ensure it does not panic when w is nil.
	s := Summarize(sampleEntries())
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSummary panicked with nil writer: %v", r)
		}
	}()
	PrintSummary(s, nil)
}
