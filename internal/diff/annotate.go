package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// Annotation holds a human-readable note attached to a specific secret path and version.
type Annotation struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	Author    string    `json:"author"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// AnnotationStore is an in-memory collection of annotations keyed by path.
type AnnotationStore struct {
	entries map[string][]Annotation
}

// NewAnnotationStore creates an empty AnnotationStore.
func NewAnnotationStore() *AnnotationStore {
	return &AnnotationStore{entries: make(map[string][]Annotation)}
}

// Add appends an annotation to the store, setting CreatedAt if zero.
func (s *AnnotationStore) Add(a Annotation) {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	s.entries[a.Path] = append(s.entries[a.Path], a)
}

// Get returns all annotations for the given path, sorted by CreatedAt ascending.
func (s *AnnotationStore) Get(path string) []Annotation {
	list := s.entries[path]
	sorted := make([]Annotation, len(list))
	copy(sorted, list)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})
	return sorted
}

// ForVersion returns annotations for a specific path and version.
func (s *AnnotationStore) ForVersion(path string, version int) []Annotation {
	var out []Annotation
	for _, a := range s.entries[path] {
		if a.Version == version {
			out = append(out, a)
		}
	}
	return out
}

// PrintAnnotations writes annotations for a path to stdout.
func PrintAnnotations(store *AnnotationStore, path string) {
	FprintAnnotations(os.Stdout, store, path)
}

// FprintAnnotations writes annotations for a path to the given writer.
func FprintAnnotations(w io.Writer, store *AnnotationStore, path string) {
	annotations := store.Get(path)
	if len(annotations) == 0 {
		fmt.Fprintf(w, "No annotations for %s\n", path)
		return
	}
	fmt.Fprintf(w, "Annotations for %s:\n", path)
	for _, a := range annotations {
		fmt.Fprintf(w, "  [v%d] %s (%s): %s\n",
			a.Version,
			a.Author,
			a.CreatedAt.Format(time.RFC3339),
			a.Note,
		)
	}
}
