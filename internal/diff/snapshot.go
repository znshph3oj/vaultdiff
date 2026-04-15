package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of a secret at a point in time.
type Snapshot struct {
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	CapturedAt time.Time        `json:"captured_at"`
	Data      map[string]string `json:"data"`
}

// NewSnapshot creates a new Snapshot for the given path, version, and data.
func NewSnapshot(path string, version int, data map[string]string) *Snapshot {
	return &Snapshot{
		Path:       path,
		Version:    version,
		CapturedAt: time.Now().UTC(),
		Data:       data,
	}
}

// SaveSnapshot writes a snapshot to a JSON file at the given filepath.
func SaveSnapshot(s *Snapshot, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from a JSON file at the given filepath.
func LoadSnapshot(filepath string) (*Snapshot, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// DiffSnapshots compares two snapshots and returns the list of changes.
func DiffSnapshots(old, new *Snapshot) []Change {
	return Compare(old.Data, new.Data)
}
