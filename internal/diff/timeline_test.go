package diff

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/vault"
)

type mockSecretReader struct {
	metadata *vault.SecretMetadata
	versions map[int]map[string]string
	metaErr  error
}

func (m *mockSecretReader) GetSecretVersion(_ context.Context, _ string, v int) (map[string]string, error) {
	data, ok := m.versions[v]
	if !ok {
		return nil, errors.New("version not found")
	}
	return data, nil
}

func (m *mockSecretReader) GetSecretMetadata(_ context.Context, _ string) (*vault.SecretMetadata, error) {
	if m.metaErr != nil {
		return nil, m.metaErr
	}
	return m.metadata, nil
}

func sampleMetadata() *vault.SecretMetadata {
	now := time.Now()
	return &vault.SecretMetadata{
		Path:           "myapp/config",
		CurrentVersion: 3,
		OldestVersion:  1,
		Versions: map[int]*vault.VersionMetadata{
			1: {Version: 1, CreatedTime: now.Add(-2 * time.Hour)},
			2: {Version: 2, CreatedTime: now.Add(-1 * time.Hour)},
			3: {Version: 3, CreatedTime: now},
		},
	}
}

func TestBuildTimeline_DetectsChanges(t *testing.T) {
	r := &mockSecretReader{
		metadata: sampleMetadata(),
		versions: map[int]map[string]string{
			1: {"key": "v1"},
			2: {"key": "v2"},
			3: {"key": "v2", "new": "added"},
		},
	}

	tl, err := BuildTimeline(context.Background(), r, "myapp/config", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tl) != 2 {
		t.Fatalf("expected 2 timeline entries, got %d", len(tl))
	}
	if len(tl[0].Changes) != 1 {
		t.Errorf("expected 1 change in entry 0, got %d", len(tl[0].Changes))
	}
	if tl[0].FromVersion != 1 || tl[0].ToVersion != 2 {
		t.Errorf("unexpected version range: %d→%d", tl[0].FromVersion, tl[0].ToVersion)
	}
}

func TestBuildTimeline_SingleVersion_ReturnsEmpty(t *testing.T) {
	meta := &vault.SecretMetadata{
		CurrentVersion: 1,
		OldestVersion:  1,
		Versions:       map[int]*vault.VersionMetadata{1: {Version: 1}},
	}
	r := &mockSecretReader{metadata: meta, versions: map[int]map[string]string{1: {"k": "v"}}}

	tl, err := BuildTimeline(context.Background(), r, "myapp/config", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tl) != 0 {
		t.Errorf("expected empty timeline, got %d entries", len(tl))
	}
}

func TestBuildTimeline_MetadataError(t *testing.T) {
	r := &mockSecretReader{metaErr: errors.New("vault unreachable")}
	_, err := BuildTimeline(context.Background(), r, "myapp/config", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPrintTimeline_ContainsVersionRange(t *testing.T) {
	tl := Timeline{
		{FromVersion: 1, ToVersion: 2, At: time.Now(), Changes: []Change{
			{Key: "password", OldValue: "old", NewValue: "new", Type: Modified},
		}},
	}
	var buf bytes.Buffer
	PrintTimeline(&buf, tl, "myapp/config")
	out := buf.String()
	if !strings.Contains(out, "v1 → v2") {
		t.Errorf("expected version range in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1 change(s)") {
		t.Errorf("expected change count in output, got:\n%s", out)
	}
}
