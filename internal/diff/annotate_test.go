package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAnnotationStore_AddAndGet(t *testing.T) {
	store := NewAnnotationStore()
	now := time.Now().UTC()
	store.Add(Annotation{Path: "secret/app", Version: 1, Author: "alice", Note: "initial", CreatedAt: now})
	store.Add(Annotation{Path: "secret/app", Version: 2, Author: "bob", Note: "updated key", CreatedAt: now.Add(time.Minute)})

	anns := store.Get("secret/app")
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	if anns[0].Author != "alice" {
		t.Errorf("expected first author alice, got %s", anns[0].Author)
	}
}

func TestAnnotationStore_GetSortedByTime(t *testing.T) {
	store := NewAnnotationStore()
	later := time.Now().UTC()
	earlier := later.Add(-time.Hour)
	store.Add(Annotation{Path: "secret/x", Version: 1, Author: "b", Note: "second", CreatedAt: later})
	store.Add(Annotation{Path: "secret/x", Version: 1, Author: "a", Note: "first", CreatedAt: earlier})

	anns := store.Get("secret/x")
	if anns[0].Author != "a" {
		t.Errorf("expected sorted order: first entry should be 'a', got %s", anns[0].Author)
	}
}

func TestAnnotationStore_ForVersion(t *testing.T) {
	store := NewAnnotationStore()
	store.Add(Annotation{Path: "secret/svc", Version: 1, Author: "alice", Note: "v1 note"})
	store.Add(Annotation{Path: "secret/svc", Version: 2, Author: "bob", Note: "v2 note"})

	v1 := store.ForVersion("secret/svc", 1)
	if len(v1) != 1 || v1[0].Note != "v1 note" {
		t.Errorf("unexpected ForVersion result: %+v", v1)
	}

	v2 := store.ForVersion("secret/svc", 2)
	if len(v2) != 1 || v2[0].Author != "bob" {
		t.Errorf("unexpected ForVersion result: %+v", v2)
	}
}

func TestAnnotationStore_SetTimestampWhenZero(t *testing.T) {
	store := NewAnnotationStore()
	store.Add(Annotation{Path: "secret/ts", Version: 1, Author: "dev", Note: "auto ts"})
	anns := store.Get("secret/ts")
	if anns[0].CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set automatically")
	}
}

func TestFprintAnnotations_NoAnnotations(t *testing.T) {
	store := NewAnnotationStore()
	var buf bytes.Buffer
	FprintAnnotations(&buf, store, "secret/missing")
	if !strings.Contains(buf.String(), "No annotations") {
		t.Errorf("expected 'No annotations' message, got: %s", buf.String())
	}
}

func TestFprintAnnotations_PrintsEntries(t *testing.T) {
	store := NewAnnotationStore()
	store.Add(Annotation{
		Path: "secret/db", Version: 3, Author: "carol",
		Note: "rotated password",
		CreatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	})
	var buf bytes.Buffer
	FprintAnnotations(&buf, store, "secret/db")
	out := buf.String()
	if !strings.Contains(out, "carol") {
		t.Errorf("expected author in output, got: %s", out)
	}
	if !strings.Contains(out, "rotated password") {
		t.Errorf("expected note in output, got: %s", out)
	}
	if !strings.Contains(out, "v3") {
		t.Errorf("expected version in output, got: %s", out)
	}
}
